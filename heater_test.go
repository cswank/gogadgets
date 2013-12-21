package gogadgets

import (
	"time"
	"bitbucket.com/cswank/gogadgets/utils"
	"testing"
)


func TestCreateHeater(t *testing.T) {
	_ = Heater{
		gpio: &FakeOutput{},
		target: 100.0,
	}
}

func getMessage(val float64) *Message {
	return &Message{
		Name: "temperature",
		Value: Value{
			Value: val,
		},
	}
}

func TestHeater(t *testing.T) {
	if !utils.FileExists("/sys/class/gpio/export") {
		return //not a beaglebone
	}
	g, err := NewGPIO(&Pin{Port:"9", Pin:"14", Direction:"out"})
	if err != nil {
		t.Fatal(err)
	}
	h := Heater{
		gpio: g,
		target: 100.0,
	}
	v := &Value{
		Value: 85.0,
		Units: "C",
	}
	h.On(v)
	time.Sleep(1 * time.Second)
	msg := getMessage(84.0)
	h.Update(msg)
	time.Sleep(5 * time.Second)
	msg = getMessage(84.5)
	h.Update(msg)
	time.Sleep(5 * time.Second)
	msg = getMessage(85.5)
	h.Update(msg)
	time.Sleep(5 * time.Second)
	msg = getMessage(82.0)
	h.Update(msg)
	time.Sleep(5 * time.Second)
	h.Off()
	
}

