package gogadgets

import (
	"errors"
)

type InputDeviceFactory func(pin *Pin) (InputDevice, error)

//Inputdevices are started as goroutines by the Gadget
//that contains it.
type InputDevice interface {
	Start(<-chan Message, chan<- Value)
	GetValue() *Value
	Config() ConfigHelper
}

type Poller interface {
	Wait() (bool, error)
}

func NewInputDevice(pin *Pin) (dev InputDevice, err error) {
	if pin.Type == "thermometer" {
		dev, err = NewThermometer(pin)
	} else if pin.Type == "switch" {
		dev, err = NewSwitch(pin)
	} else {
		err = errors.New("invalid pin type")
	}
	return dev, err
}
