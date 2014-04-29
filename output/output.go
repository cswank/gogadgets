package output

import (
	"bitbucket.org/cswank/gogadgets/models"
	"errors"
)

type OutputDeviceFactory func(pin *models.Pin) (OutputDevice, error)

//Outputdevices turn things on and off.  Currently the
//only
type OutputDevice interface {
	On(val *models.Value) error
	Off() error
	Update(msg *models.Message)
	Status() interface{}
	Config() models.Pin
}

func NewOutputDevice(pin *models.Pin) (dev OutputDevice, err error) {
	if pin.Type == "gpio" {
		dev, err = NewGPIO(pin)
	} else if pin.Type == "heater" {
		dev, err = NewHeater(pin)
	} else if pin.Type == "cooler" {
		dev, err = NewCooler(pin)
	} else if pin.Type == "recorder" {
		dev, err = NewRecorder(pin)
	} else if pin.Type == "pwm" {
		dev, err = NewPWM(pin)
	} else if pin.Type == "motor" {
		dev, err = NewMotor(pin)
	} else {
		dev, err = nil, errors.New("invalid pin type")
	}
	return dev, err
}

