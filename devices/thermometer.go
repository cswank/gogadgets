package devices

import (
	//"bitbucket.com/cswank/gogadgets/utils"
	"bitbucket.com/cswank/gogadgets/models"
)

type Thermometer struct {
	InputDevice
	units string
}

func NewThermometer(pin *models.Pin) (*GPIO, error) {
	return nil, nil
}

func (t *Thermometer) Start(out <-chan models.Message) {
	
}
