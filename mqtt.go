package main

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var BoolToMQTT = map[bool]string{
	true:  "ON",
	false: "OFF",
}

const (
	TOPIC = "finn/meeting_active/state"
)

type Mqtt struct {
	State  bool
	client mqtt.Client
}

func NewMqtt(broker, port, username, password string) *Mqtt {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", broker, port))
	opts.SetClientID("go_mqtt_client")
	opts.SetUsername(username)
	opts.SetPassword(password)
	return &Mqtt{
		client: mqtt.NewClient(opts),
		State:  false,
	}
}

type State struct {
	State string `json:"state"`
}

func (m *Mqtt) setState(newState bool) {
	fmt.Printf("Changing state to %v", newState)
	m.State = newState
	token := m.client.Publish(TOPIC, 0, true, BoolToMQTT[m.State])
	token.Wait()
}
