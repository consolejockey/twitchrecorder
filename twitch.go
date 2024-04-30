package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"
)

const (
	TOKEN_API  = "https://id.twitch.tv/oauth2/token"
	STREAM_API = "https://api.twitch.tv/helix/streams?user_login="
)

type TwitchClient struct {
	ClientID     string
	ClientSecret string
	AccessToken  string
	client       *http.Client
}

func NewTwitch(clientID, clientSecret string) (*TwitchClient, error) {
	t := &TwitchClient{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		client:       &http.Client{},
	}
	t.AccessToken = t.RefreshAccessToken()
	return t, nil
}

func (t TwitchClient) PrintClientInfo() {
	log.Println("Twitch Info:")
	log.Printf("ClientID: %s\n", t.ClientID)
	log.Printf("ClientSecret: %s\n", t.ClientSecret)
	log.Printf("AccessToken: %s\n", t.AccessToken)
}

func (t TwitchClient) generateHeaders() map[string]string {
	headers := map[string]string{
		"Client-ID":     t.ClientID,
		"Authorization": "Bearer " + t.AccessToken,
	}
	return headers
}

func (t *TwitchClient) RefreshAccessToken() string {
	body := url.Values{}
	body.Set("client_id", t.ClientID)
	body.Set("client_secret", t.ClientSecret)
	body.Set("grant_type", "client_credentials")

	resp, err := t.client.PostForm(TOKEN_API, body)
	if err != nil {
		log.Printf("Could not request a new access token! Error: %s\n", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading access token response:", err)
	}

	var responseMap map[string]interface{}
	if err := json.Unmarshal(respBody, &responseMap); err != nil {
		log.Fatal("Error parsing access token response:", err)
	}

	accessToken, ok := responseMap["access_token"].(string)
	if !ok {
		log.Println("Access token not found in response!")
	}

	t.AccessToken = accessToken
	if accessToken != "" {
		log.Println("Access token refreshed:", t.AccessToken)
	} else {
		log.Println("Error! Access token is empty string.")
	}

	return accessToken
}

func (t *TwitchClient) IsLive(streamerName string) bool {
	var headers = t.generateHeaders()

	req, err := http.NewRequest("GET", STREAM_API+streamerName, nil)
	if err != nil {
		log.Println("Error creating request to check live status:", err)
		return false
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := t.client.Do(req)
	if err != nil {
		log.Println("Error sending request for live status:", err)
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		log.Println("Status 401: Unauthorized")
		t.RefreshAccessToken()
		return false
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response for live status:", err)
		return false
	}

	var responseData map[string]interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		log.Println("Error parsing live status response:", err)
		return false
	}

	data, ok := responseData["data"].([]interface{})
	if !ok {
		log.Println("Error parsing data field of live status response:", err)
		return false
	}

	return len(data) > 0
}
