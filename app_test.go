package gogadgets

import (
	"fmt"
	"time"
	"testing"
)

type FakeOutput struct {
	OutputDevice
	on bool
}

func (f *FakeOutput) Update(msg *Message) {
	
}

func (f *FakeOutput) On(val *Value) error {
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
	Poller
	val bool
}

func (f *FakePoller) Wait() (bool, error) {
	time.Sleep(100 * time.Millisecond)
	f.val = !f.val
	return f.val, nil
}

func TestGetGadgets(t *testing.T) {
	// if !utils.FileExists("/sys/class/gpio/export") {
	// 	return //not a beaglebone
	// }
	// configs := []*Config{
	// 	&Config{
	// 		Type: "gpio",
	// 		Location: "tank",
	// 		Name: "pump",
	// 		Pin: Pin{
	// 			Port: "9",
	// 			Pin: "15",
	// 		},
	// 	},
	// 	&Config{
	// 		Type: "switch",
	// 		Location: "tank",
	// 		Name: "switch",
	// 		Pin: Pin{
	// 			Port: "9",
	// 			Pin: "16",
	// 			Value: "7.5",
	// 			Units: "liters",
	// 		},
	// 	},
	// }
}

func TestGadgets(t *testing.T) {
	p := &Gadget{
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
	s := &Gadget{
		Location: location,
		Name: name,
		Input: &Switch{
			GPIO: poller,
			Value: 5.0,
			Units: "liters",
		},
		UID: fmt.Sprintf("%s %s", location, name),
	}
	a := App{
		gadgets: []GoGadget{p, s},
	}
	stop := make(chan bool)
	go a.Start(stop)
	time.Sleep(1 * time.Second)
}
