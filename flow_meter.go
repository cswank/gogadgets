package gogadgets

import (
	"errors"
	"fmt"
	"log"
	"time"
)

//Flow meter waits for a high pulse from a gpio pin
//then caclulates the flow based on the time between
//high pulses.
type FlowMeter struct {
	GPIO Poller
	//Value represents the total volume represented by
	//a pulse (rate is calculated from this total volume).
	Value float64
	value float64
	ts    time.Time
	Units string
	out   chan<- Value
}

func NewFlowMeter(pin *Pin) (InputDevice, error) {
	pin.Direction = "in"
	pin.Value = true
	gpio, err := NewGPIO(pin)
	if err != nil {
		return nil, err
	}
	poller, ok := gpio.(Poller)
	if !ok {
		return nil, errors.New(fmt.Sprintf("couldn't create a poller: %s", pin))
	}

	return &FlowMeter{
		GPIO:  poller,
		Value: pin.Value.(float64),
		Units: pin.Units,
	}, nil
}

func (f *FlowMeter) Config() ConfigHelper {
	return ConfigHelper{}
}

func (f *FlowMeter) wait(err chan<- error) {
	for {
		v, e := f.GPIO.Wait()
		if !v {
			continue
		}
		if e != nil {
			err <- e
			continue
		}
		t1 := f.ts
		f.ts = time.Now()
		if t1.Year() == 1 {
			continue
		}
		f.value = f.Value / float64(f.ts.Sub(t1).Seconds())
		f.out <- Value{
			Value: f.value,
			Units: f.Units,
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (f *FlowMeter) SendValue() {
	f.out <- *f.GetValue()
}

func (f *FlowMeter) GetValue() *Value {
	return &Value{
		Value: f.value,
		Units: f.Units,
	}
}

func (f *FlowMeter) Start(in <-chan Message, out chan<- Value) {
	f.out = out
	err := make(chan error)
	f.SendValue()
	keepGoing := true
	go f.wait(err)
	for keepGoing {
		select {
		case <-in:
			//do nothing
		case e := <-err:
			log.Println(e)
		}
	}
}
