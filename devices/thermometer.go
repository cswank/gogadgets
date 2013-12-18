package devices

import (
	//"bitbucket.com/cswank/gogadgets/utils"
	"bitbucket.com/cswank/gogadgets/models"
)

type Thermometer struct {
	InputDevice
	units string
}

func NewThermometer(pin *models.Pin) (*Thermometer, error) {
	return nil, nil
}

func (t *Thermometer) Start(stop <-chan bool, out chan<- models.Value) {
	
}
