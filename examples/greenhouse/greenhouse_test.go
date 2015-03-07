package main

import (
	"testing"
	"time"

	"github.com/cswank/gogadgets"
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
		Type:     "update",
		Location: "greenhouse",
		Name:     "temperature",
		Value: gogadgets.Value{
			Value: 14.0,
			Units: "C",
		},
	}
	out <- msg

	msg = <-in
	if msg.Body != "turn on bed 1 pump" {
		t.Error("pump should be on", msg.Body)
	}
	msg = <-in
	if msg.Body != "turn on bed 2 pump" {
		t.Error("pump should be on", msg.Body)
	}
	msg = <-in
	if msg.Body != "turn on bed 3 pump" {
		t.Error("pump should be on", msg.Body)
	}

	msg = gogadgets.Message{
		Type:     "update",
		Location: "bed 1",
		Name:     "switch",
		Value: gogadgets.Value{
			Value: 0.0,
		},
	}
	go func() {
		time.Sleep(10 * time.Millisecond)
		out <- msg
	}()
	msg = <-in
	if msg.Body != "turn off bed 1 pump" {
		t.Error("pump should be off", msg.Body)
	}
	msg = <-in
	if msg.Body != "turn on bed 1 pump" {
		t.Error("pump should be on", msg.Body)
	}
}
