package main

import (
	"fmt"
	"time"
	"bitbucket.com/cswank/gogadgets/gadgets"
	"bitbucket.com/cswank/gogadgets/models"
	"bitbucket.com/cswank/gogadgets/devices"
	"testing"
)

type FakeOutput struct {
	devices.OutputDevice
	on bool
}

func (f *FakeOutput) Update(msg *models.Message) {
	
}

func (f *FakeOutput) On(val *models.Value) error {
	f.on = true
	return nil
}

func (f *FakeOutput) Off() error {
	f.on = false
	return nil
}

func (f *FakeOutput) Status() interface{} {
	return f.on
}

type FakePoller struct {
	devices.Poller
	val bool
}

func (f *FakePoller) Wait() (bool, error) {
	time.Sleep(100 * time.Millisecond)
	f.val = !f.val
	return f.val, nil
}

func TestGadgets(t *testing.T) {
	p := &gadgets.Gadget{
		Location: "tank",
		Name: "pump",
		OnCommand: fmt.Sprintf("turn on %s %s", "tank", "pump"),
		OffCommand: fmt.Sprintf("turn off %s %s", "tank", "pump"),
		Output: &FakeOutput{},
		UID: fmt.Sprintf("%s %s", "tank", "pump"),
	}
	location := "tank"
	name := "switch"
	poller := &FakePoller{}
	s := &gadgets.Gadget{
		Location: location,
		Name: name,
		Input: &devices.Switch{
			GPIO: poller,
			Value: 5.0,
			Units: "liters",
		},
		UID: fmt.Sprintf("%s %s", location, name),
	}
	a := App{
		gadgets: []models.Gadget{p, s},
	}
	stop := make(chan bool)
	go a.Start(stop)
	time.Sleep(1 * time.Second)
}
