package output

import (
	"fmt"
	"testing"
	"time"
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
	}
	input := make(chan gogadgets.Message)
	output := make(chan gogadgets.Message)
	go g.Start(input, output)
	msg := gogadgets.Message{
		Type: "command",
		Command: "turn on lab led",
	}
	input<- msg
	for !g.Output.Status() {
		fmt.Println("still off")
		time.Sleep(10 * time.Millisecond)
	}
	msg = gogadgets.Message{
		Type: "command",
		Command: "shutdown",
	}
	input<- msg
	for g.Output.Status() {
		//device should turn off when shutdown message is received
		fmt.Println("still on")
		time.Sleep(10 * time.Millisecond)
	}
}
