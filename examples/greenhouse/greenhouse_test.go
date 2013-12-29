package main

import (
	"time"
	"testing"
	"bitbucket.com/cswank/gogadgets"
)

func TestGreenhouse(t *testing.T) {
	sleepTimes := map[string]time.Duration{
		"bed 1": 1,
		"bed 2": 1,
		"bed 3": 1,
	}
	g := Greenhouse{sleepTimes: sleepTimes}
	in := make(chan gogadgets.Message)
	out := make(chan gogadgets.Message)
	go g.Start(out, in)
	msg := gogadgets.Message{
		Type: "update",
		Location: "greenhouse",
		Name: "temperature",
		Value: gogadgets.Value{
			Value: 14.0,
			Units: "C",
		},
	}
	out<- msg

	msg = gogadgets.Message{
		Type: "update",
		Location: "bed 1",
		Name: "switch",
		Value: gogadgets.Value{
			Value: false,
		},
	}

	out<- msg
	msg = <-in
	if msg.Body != "turn off bed 1 pump" {
		t.Error("pump should be off", msg.Body)
	}
	msg = <-in
	if msg.Body != "turn on bed 1 pump" {
		t.Error("pump should be off", msg.Body)
	}
	
}
