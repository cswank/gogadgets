package gogadgets

import (
	"fmt"
	"errors"
	"log"
	"time"
)

//Switch is an input device that waits for a GPIO pin
//to change value (1 to 0 or 0 to 1).  When that change
//happens it sends an update to the rest of the system.
type Switch struct {
	GPIO  Poller
	Value interface{}
	Units string
	out   chan<- Value
}

func NewSwitch(pin *Pin) (InputDevice, error) {
	var err error
	var s *Switch
	gpio, err := NewGPIO(pin)
	poller, ok := gpio.(Poller)
	if !ok {
		return nil, errors.New(fmt.Sprintf("couldn't create a poller: %s", pin))
	}
	if err == nil {
		s = &Switch{
			GPIO:  poller,
			Value: pin.Value,
			Units: pin.Units,
		}
	}
	return s, err
}

//The GPIO does the real waiting here.  This wraps it and adds
//a delay so that the inevitable bounce in the signal from the
//physical device is ignored.
func (s *Switch) wait(out chan<- interface{}, err chan<- error) {
	val, e := s.GPIO.Wait()
	if e != nil {
		err <- e
	} else {
		if val {
			out <- s.Value
		} else {
			if s.Value == true {
				out <- false
			} else {
				out <- 0.0
			}
		}
	}
	time.Sleep(200 * time.Millisecond)
}

func (s *Switch) SendValue() {
	s.out <- Value{
		Value: s.Value,
		Units: s.Units,
	}
}

func (s *Switch) GetValue() *Value {
	return &Value{
		Value: s.Value,
		Units: s.Units,
	}
}

func (s *Switch) Start(in <-chan Message, out chan<- Value) {
	s.out = out
	value := make(chan interface{})
	err := make(chan error)
	keepGoing := true
	for keepGoing {
		go s.wait(value, err)
		select {
		case <-in:
			//do nothing
		case val := <-value:
			s.Value = val
			s.SendValue()
		case e := <-err:
			log.Println(e)
		}
	}
}
