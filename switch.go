package gogadgets

import (
	"log"
	"time"
)

//Switch is an input device that waits for a GPIO pin
//to change value (1 to 0 or 0 to 1).  When that change
//happens it sends an update to the rest of the system.
type Switch struct {
	GPIO  Poller
	Value bool
	Units string
	out   chan<- Value

	delay    time.Duration
	lastRead time.Time
}

func NewSwitch(pin *Pin, opts ...func(InputDevice) error) (InputDevice, error) {
	pin.Direction = "in"
	pin.Edge = "both"
	var err error
	gpio, err := newGPIO(pin)
	if err != nil {
		return nil, err
	}

	return &Switch{
		GPIO:  gpio,
		Units: pin.Units,
		delay: time.Duration(50 * time.Millisecond),
	}, nil
}

//The GPIO does the real waiting here.  This wraps it and adds
//a delay so that the inevitable bounce in the signal from the
//physical device is ignored.
func (s *Switch) wait(out chan<- bool) {
	for {
		e := s.doWait()
		if e != nil {
			log.Printf("gpio wait error: %s", e)
		} else {
			out <- true
		}
	}
}

func (s *Switch) doWait() error {
	n := time.Now()
	if n.Sub(s.lastRead) < s.delay {
		time.Sleep(s.delay)
		return s.doWait()
	}

	err := s.GPIO.Wait()
	s.lastRead = n
	return err
}

func (s *Switch) readValue() {
	v := s.GPIO.Status()
	s.Value = v["gpio"]
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
	value := make(chan bool)
	s.readValue()
	s.SendValue()
	keepGoing := true
	go s.wait(value)
	for keepGoing {
		select {
		case <-in:
			//do nothing
		case val := <-value:
			s.Value = val
			s.readValue()
			s.SendValue()
		}
	}
}
