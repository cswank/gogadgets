package devices

import (
	"errors"
)

type OutputDevice interface {
	On() error
	Off() error
	Status() bool
}

type InputDevice interface {
	Status() bool
}

func NewOutputDevice(pin *Pin) (OutputDevice, error) {
	if pin.Type == "gpio" {
		return NewGPIO(pin)
	} else {
		return nil, errors.New("invalid pin type")
	}
}

