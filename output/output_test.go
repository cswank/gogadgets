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
	msg := gogadgets.Message{
		Type: "command",
		Command: "shutdown",
	}
	input := make(chan gogadgets.Message)
	output := make(chan gogadgets.Message)
	go g.Start(input, output)
	input<- msg
}
