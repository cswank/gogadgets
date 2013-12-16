package devices

import (
	"errors"
	"bitbucket.com/cswank/gogadgets"
)

type OutputDevice interface {
	On(val *gogadgets.Value) error
	Off() error
	Update(msg *gogadgets.Message)
	Status() bool
}

type InputDevice interface {
	Status() bool
	Wait(value interface{})
}

func NewOutputDevice(pin *Pin) (OutputDevice, error) {
	if pin.Type == "gpio" {
		return NewGPIO(pin)
	} else {
		return nil, errors.New("invalid pin type")
	}
}

