package devices

import (
	"bitbucket.com/cswank/gogadgets/pins"
	"errors"
)

type OutputDevice interface {
	On() error
	Off() error
}

func NewOutputDevice(pin *pins.Pin) (OutputDevice, error) {
	if pin.Type == "gpio" {
		return NewGPOutput(pin)
	} else {
		return nil, errors.New("invalid pin type")
	}
}

