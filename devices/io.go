package devices

import (
	"errors"
	"bitbucket.com/cswank/gogadgets/models"
)

type OutputDevice interface {
	On(val *models.Value) error
	Off() error
	Update(msg *models.Message)
	Status() interface{}
}

type InputDevice interface {
	Start(<-chan bool, chan<- models.Value)
}

type Poller interface {
	Wait() (bool, error)
}

func NewOutputDevice(pin *models.Pin) (dev OutputDevice, err error) {
	if pin.Type == "gpio" {
		dev, err = NewGPIO(pin)
	} else if pin.Type == "heater" {
		dev, err = NewHeater(pin)
	} else {
		dev, err = nil, errors.New("invalid pin type")
	}
	return dev, err
}

func NewInputDevice(pin *models.Pin) (dev InputDevice, err error) {
	if pin.Type == "thermometer" {
		dev, err = NewThermometer(pin)
	} else if pin.Type == "swtich" {
		dev, err =  NewSwitch(pin)
	} else {
		dev, err = nil, errors.New("invalid pin type")
	}
	return dev, err
}

