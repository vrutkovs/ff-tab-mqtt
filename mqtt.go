package main

var BoolToMQTT = map[bool]string{
	true:  "on",
	false: "off",
}

const (
	TOPIC = "finn/meeting_active"
)

type Mqtt struct {
	State bool
}

func (m *Mqtt) setState(newState bool) {
	m.State = newState
}
