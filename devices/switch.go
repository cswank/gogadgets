package devices

import (
	"syscall"
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


func (s *Switch) Start(in <-chan models.Message, out chan-> models.Message) {
	
}

func FD_SET(fd int, p *syscall.FdSet) {
        p.Bits[fd/32] |= 1 << (uint(fd) % 32)
}
