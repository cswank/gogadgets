package devices

import (
	"log"
	"bitbucket.com/cswank/gogadgets/models"
)

type Switch struct {
	InputDevice
	GPIO Poller
	Value float32
	Units string
}

func NewSwitch(pin *models.Pin) (s *Switch, err error) {
	gpio, err := NewGPIO(pin)
	if err == nil {
		s = &Switch{
			GPIO:gpio,
			Value: pin.Value.(float32),
			Units: pin.Units,
		}
	}
	return s, err
}

func (s *Switch) wait(out chan<- float32, err chan<- error) {
	val, e := s.GPIO.Wait()
	if e != nil {
		err<- e
	} else {
		if val {
			out<- s.Value
		} else {
			out<- 0.0
		}
	}
}

func (s *Switch) Start(stop <-chan bool, out chan<- models.Value) {
	value := make(chan float32)
	err := make(chan error)
	for {
		go s.wait(value, err)
		select {
		case <-stop:
			return
		case val := <-value:
			out<- models.Value{
				Value: val,
				Units: s.Units,
			}
		case e := <-err:
			log.Println(e)
		}
	}
}
