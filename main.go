package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var meetingUrls = []string{
	"https://meet.google.com",
	"https://bluejeans.com",
}

func main() {
	// Read env vars
	mqttHost, ok := os.LookupEnv("MQTT_HOST")
	if !ok {
		panic("MQTT_HOST unset")
	}
	mqttPort, ok := os.LookupEnv("MQTT_PORT")
	if !ok {
		panic("MQTT_PORT unset")
	}
	mqttUsername, ok := os.LookupEnv("MQTT_USER")
	if !ok {
		panic("MQTT_USER unset")
	}
	mqttPassword, ok := os.LookupEnv("MQTT_PASSWORD")
	if !ok {
		panic("MQTT_PASSWORD unset")
	}
	profileDir, ok := os.LookupEnv("FF_PROFILE")
	if !ok {
		panic("FF_PROFILE unset")
	}

	// Make an mqtt client
	mqtt := NewMqtt(mqttHost, mqttPort, mqttUsername, mqttPassword)
	if token := mqtt.client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	// Check for meeting every 5 seconds
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})

	// Close on SIGTERM
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		close(quit)
		os.Exit(1)
	}()

	// Every time ticker sends a channel message run actions
	fmt.Printf("Started main loop")
	for {
		select {
		case <-ticker.C:
			meetingFound := meetingSensor(profileDir)
			if mqtt.State != meetingFound {
				mqtt.setState(meetingFound)
			}
		case <-quit:
			mqtt.client.Disconnect(250)
			ticker.Stop()
			return
		}
	}
}

func meetingSensor(profileDir string) bool {
	urls, err := collectUrls(profileDir)
	if err != nil {
		fmt.Printf("Error fetching urls from FF: %v", err)
	}
	found := false
	for m := range meetingUrls {
		for i := range urls {
			if strings.HasPrefix(urls[i], meetingUrls[m]) {
				found = true
				break
			}
		}
	}
	return found
}
