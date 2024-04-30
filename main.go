package main

import (
	"log"
	"path/filepath"
	"slices"
	"time"
)

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
