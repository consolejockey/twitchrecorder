package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
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
