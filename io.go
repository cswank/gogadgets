package gogadgets

import (
	"errors"
)

type OutputDevice interface {
	On(val *Value) error
	Off() error
	Update(msg *Message)
	Status() interface{}
}

type InputDevice interface {
	Start(<-chan bool, chan<- Value)
}

type Poller interface {
	Wait() (bool, error)
}

func NewOutputDevice(pin *Pin) (dev OutputDevice, err error) {
	if pin.Type == "gpio" {
		dev, err = NewGPIO(pin)
	} else if pin.Type == "heater" {
		dev, err = NewHeater(pin)
	} else {
		dev, err = nil, errors.New("invalid pin type")
	}
	return dev, err
}

func NewInputDevice(pin *Pin) (dev InputDevice, err error) {
	if pin.Type == "thermometer" {
		dev, err = NewThermometer(pin)
	} else if pin.Type == "switch" {
		dev, err =  NewSwitch(pin)
	} else {
		err = errors.New("invalid pin type")
	}
	return dev, err
}

