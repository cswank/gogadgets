package output

import (
	"time"
	"bitbucket.com/cswank/gogadgets"
	"testing"
)

func sendVolumeMessage(out chan<- gogadgets.Message, val float64) {
	time.Sleep(10 * time.Millisecond)
	l := gogadgets.Location{
		Input: map[string]gogadgets.Device{
			"volume": gogadgets.Device{
				Units: "liters",
				Value: val,
			},
		},
	}
	msg := gogadgets.Message{
		Sender: "tank volume",
		Type: gogadgets.STATUS,
		Locations: map[string]gogadgets.Location{"tank": l},
	}
	out<- msg
}

func TestVolume(t *testing.T) {
	tr := Trigger{
		location: "tank",
		name: "valve",
		operator: ">=",
		command: " to 5 liters",
		offCommand: "stop filling tank",
	}
	in := make(chan gogadgets.Message)
	out := make(chan gogadgets.Message)
	go tr.Start(in, out)
	go sendVolumeMessage(out, 5.0)
	msg := <-in
	if msg.Body != "stop filling tank" {
		t.Error(msg)
	}
}

func TestStripCommand(t *testing.T) {
	tr := Trigger{
		location: "tank",
		name: "valve",
		operator: ">=",
		command: " to 5 liters",
		offCommand: "stop filling tank",
	}
	tr.stripCommand()
	if tr.command != "5 liters" {
		t.Error(tr.command)
	}
}

func TestGetValue(t *testing.T) {
	tr := Trigger{
		location: "tank",
		name: "valve",
		operator: ">=",
		command: " to 5 liters",
		offCommand: "stop filling tank",
	}
	val, unit, err := tr.getValue()
	if err != nil {
		t.Error(err)
	}
	if val != 5.0 {
		t.Error("incorrect value", val)
	}
	if unit != "liters" {
		t.Error("incorrect unit", unit)
	}
}

func TestGetTimeValue(t *testing.T) {
	tr := Trigger{
		location: "tank",
		name: "valve",
		operator: ">=",
		command: "for 1.1 minutes",
		offCommand: "stop filling tank",
	}
	val, unit, err := tr.getValue()
	if err != nil {
		t.Error(err)
	}
	if val != 1.1 {
		t.Error("incorrect value", val)
	}
	if unit != "minutes" {
		t.Error("incorrect unit", unit)
	}
}
	

