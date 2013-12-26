package gogadgets

import (
	"log"
)

type Switch struct {
	InputDevice
	GPIO Poller
	Value float64
	Units string
	value float64
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

func (s *Switch) GetValue() *Value {
	return &Value{
		Value: s.value,
		Units: s.Units,
	}
}


func (s *Switch) Start(in <-chan Message, out chan<- Value) {
	value := make(chan float64)
	err := make(chan error)
	for {
		go s.wait(value, err)
		select {
		case msg := <- in:
			if msg.Type == "command" && msg.Body == "shutdown" {
				return
			} else if msg.Type == "command" && msg.Body == "status" {
				out<- Value{
					Value: s.value,
					Units: s.Units,
				}
			}
		case val := <-value:
			s.value = val
			out<- Value{
				Value: val,
				Units: s.Units,
			}
		case e := <-err:
			log.Println(e)
		}
	}
}
