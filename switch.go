package gogadgets

import (
	"log"
)

type Switch struct {
	InputDevice
	GPIO Poller
	Value float64
	Units string
}

func NewSwitch(pin *Pin) (s *Switch, err error) {
	gpio, err := NewGPIO(pin)
	if err == nil {
		s = &Switch{
			GPIO:gpio,
			Value: pin.Value.(float64),
			Units: pin.Units,
		}
	}
	return s, err
}

func (s *Switch) wait(out chan<- float64, err chan<- error) {
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

func (s *Switch) Start(stop <-chan bool, out chan<- Value) {
	value := make(chan float64)
	err := make(chan error)
	for {
		go s.wait(value, err)
		select {
		case <-stop:
			return
		case val := <-value:
			out<- Value{
				Value: val,
				Units: s.Units,
			}
		case e := <-err:
			log.Println(e)
		}
	}
}
