package output

import (
	"fmt"
	"testing"
	//"bitbucket.com/cswank/gogadgets/pins"
	"bitbucket.com/cswank/gogadgets/devices"
	"bitbucket.com/cswank/gogadgets"
)

type FakeOutput struct {
	devices.OutputDevice
	on bool
}

func (f *FakeOutput) On() error {
	f.on = true
	return nil
}

func (f *FakeOutput) Off() error {
	f.on = false
	return nil
}

func (f *FakeOutput) Status() bool {
	return f.on
}


func TestStart(t *testing.T) {
	location := "lab"
	name := "led"
	g := OutputGadget{
		Location: location,
		Name: name,
		OnCommand: fmt.Sprintf("turn on %s %s", location, name),
		OffCommand: fmt.Sprintf("turn off %s %s", location, name),
		Output: &FakeOutput{},
		uid: fmt.Sprintf("%s %s", location, name),
	}
	input := make(chan gogadgets.Message)
	output := make(chan gogadgets.Message)
	go g.Start(input, output)
	msg := gogadgets.Message{
		Type: "command",
		Body: "turn on lab led",
	}
	input<- msg
	status := <-output
	if status.Locations["lab"].Output["led"].Value != true {
		t.Error("shoulda been on", status)
	}
	
	msg = gogadgets.Message{
		Type: "command",
		Body: "shutdown",
	}
	input<- msg
	status = <-output
	if status.Locations["lab"].Output["led"].Value != false {
		t.Error("shoulda been off", status)
	}
}
