package devices

import (
	"bitbucket.com/cswank/gogadgets"
	"testing"
)

type FakeOutput struct {
	OutputDevice
	on bool
}

func (f *FakeOutput) Update(msg *gogadgets.Message) {
	
}

func (f *FakeOutput) On(val *gogadgets.Value) error {
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


func TestCreateHeater(t *testing.T) {
	_ = Heater{
		gpio: &FakeOutput{},
		target: 100.0,
	}
}

func TestHeater(t *testing.T) {
	g, err := NewGPIO(&Pin{Port:"9", Pin:"15", Direction:"out"})
	if err != nil {
		t.Fatal(err)
	}
	h := Heater{
		gpio: g,
		target: 100.0,
	}
	v := &gogadgets.Value{
		Value: 85.0,
		Units: "C",
	}
	h.On(v)
	time.Sleep(1 * time.Second())
	h.Off()
}

