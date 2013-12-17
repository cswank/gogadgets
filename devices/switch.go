package devices

import (
	//"bitbucket.com/cswank/gogadgets/utils"
	"bitbucket.com/cswank/gogadgets/models"
)

type Switch struct {
	InputDevice
	units string
}

func NewSwitch(pin *models.Pin) (*GPIO, error) {
	return nil, nil
}

func (s *Switch) Start(out <-chan models.Message) {
	
}

