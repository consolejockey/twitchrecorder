package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Recorder struct {
	Process *os.Process
}

func NewRecorder() *Recorder {
	return &Recorder{}
}

func (r *Recorder) StartRecording(streamerName, filePath, preferredQuality string) error {
	timestamp := time.Now().Format("02-01-2006_15-04-05")
	outputFile := filepath.Join(filePath, fmt.Sprintf("%s.mp4", timestamp))
	streamlinkCommand := []string{
		"streamlink",
		"-o", outputFile,
		"twitch.tv/" + streamerName,
		preferredQuality,
		"--twitch-disable-hosting",
		"--twitch-disable-ads",
		"--twitch-disable-reruns",
	}

	commandString := strings.Join(streamlinkCommand, " ")
	log.Println("Executing command:", commandString)
	cmd := exec.Command(streamlinkCommand[0], streamlinkCommand[1:]...)
	err := cmd.Start()
	// FIX: If preferred quality is not available, it still sends the command and does not return an error.
	if err != nil {
		log.Fatal("Failed to start streamlink cmd. Error:", err)
		return err
	}

	r.Process = cmd.Process
	log.Println("Recording started. Output file name is:", timestamp+".mp4")
	return nil
}

func (r *Recorder) StopRecording() error {
	if r.Process == nil {
		return nil
	}

	_, err := r.Process.Wait()
	if err != nil {
		return fmt.Errorf("error waiting for process to exit: %v", err)
	}

	r.Process = nil
	log.Println("Streamlink process terminated successfully.")
	return nil
}
