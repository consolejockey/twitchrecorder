package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"
)

type Config struct {
	ClientID         string `json:"client_id"`
	ClientSecret     string `json:"client_secret"`
	DownloadFolder   string `json:"download_folder"`
	PreferredQuality string `json:"quality"`
	Streamer         string `json:"streamer"`
}

func (config *Config) readConfig() error {

	file, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Failed to open config.json:", err)
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		log.Fatal("Failed to decode config.json:", err)
		return err
	}

	return nil
}

func (config *Config) integrityCheck() error {
	var missingFields []string

	if config.ClientID == "" {
		missingFields = append(missingFields, "ClientID")
	}
	if config.ClientSecret == "" {
		missingFields = append(missingFields, "ClientSecret")
	}
	if config.DownloadFolder == "" {
		missingFields = append(missingFields, "DownloadFolder")
	}
	if config.PreferredQuality == "" {
		missingFields = append(missingFields, "PreferredQuality")
	}
	if config.Streamer == "" {
		missingFields = append(missingFields, "Streamer")
	}

	if len(missingFields) > 0 {
		missingFieldsStr := strings.Join(missingFields, ", ")
		return fmt.Errorf("the following field(s) are missing in the configuration: %s", missingFieldsStr)
	}

	return nil
}

func main() {
	const interval = 15
	var isRecording = false

	config := Config{}
	recorder := NewRecorder()

	err := config.readConfig()
	if err != nil {
		log.Fatal("Error reading config:", err)
	}

	err = config.integrityCheck()
	if err != nil {
		log.Fatalf("Configuration is incomplete. %v", err)
	}

	downloadFolder, err := filepath.Abs(config.DownloadFolder)
	if err != nil {
		log.Fatal("Error getting absolute path for download folder:", err)
	}

	twitchClient, err := NewTwitch(config.ClientID, config.ClientSecret)
	if err != nil {
		log.Fatal("Failed to create Twitch client:", err)
	}

	twitchClient.PrintClientInfo()

	for {
		var isLive = twitchClient.IsLive(config.Streamer)

		if isLive && !isRecording {
			log.Printf("%s is now live!", config.Streamer)
			var quality = config.PreferredQuality
			var availableQualities, err = twitchClient.GetAvailableStreams(config.Streamer)
			if err != nil {
				quality = "best"
				log.Println("Could not retrieve available qualities. Defaulting to 'best'. Error:", err)
			} else if !slices.Contains(availableQualities, quality) {
				quality = "best"
				log.Printf("Could not find preferred quality '%s'. Defaulting to 'best'.", config.PreferredQuality)
			}
			recorder.StartRecording(config.Streamer, downloadFolder, quality)
			isRecording = true
		} else if !isLive && isRecording {
			log.Printf("%s has gone offline!", config.Streamer)
			err = recorder.StopRecording()
			if err != nil {
				log.Printf("Error terminating streamlink recording process. %v", err)
			}
			isRecording = false
		}

		time.Sleep(interval * time.Second)
	}
}
