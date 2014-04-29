package input

import (
	"errors"
	"bitbucket.org/cswank/gogadgets/models"
)

type InputDeviceFactory func(pin *models.Pin) (InputDevice, error)

//Inputdevices are started as goroutines by the Gadget
//that contains it.
type InputDevice interface {
	Start(<-chan models.Message, chan<- models.Value)
	GetValue() *models.Value
	Config() models.Pin
}

type Poller interface {
	Wait() (bool, error)
}

func NewInputDevice(pin *models.Pin) (dev InputDevice, err error) {
	if pin.Type == "thermometer" {
		dev, err = NewThermometer(pin)
	} else if pin.Type == "switch" {
		dev, err = NewSwitch(pin)
	} else {
		err = errors.New("invalid pin type")
	}
	return dev, err
}
