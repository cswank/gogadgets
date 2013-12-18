package devices

import (
	"log"
	"bitbucket.com/cswank/gogadgets/models"
)

type Switch struct {
	InputDevice
	gpio Poller
	units string
}

func NewSwitch(pin *models.Pin) (s *Switch, err error) {
	gpio, err := NewGPIO(pin)
	if err == nil {
		s = &Switch{gpio:gpio}
	}
	return s, err
}

func (s *Switch) wait(out chan<- bool, err chan<- error) {
	val, e := s.gpio.Wait()
	if e != nil {
		err<- e
	} else {
		out<- val
	}
}

func (s *Switch) Start(stop <-chan bool, out chan<- models.Value) {
	change := make(chan bool)
	err := make(chan error)
	for {
		go s.wait(change, err)
		select {
		case <-stop:
			return
		case val := <-change:
			out<- models.Value{
				Value: val,
				Units: s.units,
			}
		case e := <-err:
			log.Println(e)
		}
	}
}
