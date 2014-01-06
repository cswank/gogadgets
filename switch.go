package gogadgets

import (
	"log"
)

type Switch struct {
	InputDevice
	GPIO Poller
	Value float64
	Units string
	out chan<- Value
}

func NewSwitch(pin *Pin) (InputDevice, error) {
	var err error
	var s *Switch
	gpio, err := NewGPIO(pin)
	if err == nil {
		s = &Switch{
			GPIO:gpio.(Poller),
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

func (s *Switch) SendValue() {
	s.out<- Value{
		Value: s.Value,
		Units: s.Units,
	}
}

func (s *Switch) Start(in <-chan Message, out chan<- Value) {
	s.out = out
	value := make(chan float64)
	err := make(chan error)
	keepGoing := true
	for keepGoing {
		go s.wait(value, err)
		select {
		case <- in:
			//do nothing
		case val := <-value:
			s.Value = val
			s.SendValue()
		case e := <-err:
			log.Println(e)
		}
	}
}
