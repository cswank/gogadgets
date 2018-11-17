package gogadgets

import (
	"log"
	"time"
)

//FlowMeter waits for a high pulse from a gpio pin
//then caclulates the flow based on the time between
//high pulses.
type FlowMeter struct {
	GPIO Poller

	//Value represents the total volume represented by
	//a pulse (rate is calculated from this total volume).
	Value float64
	value float64
	ts    time.Time

	//min_span is the minimum number of seconds between
	//2 pulses from the physical flow meter.  It is used
	//for de-bouncing the signal.
	MinSpan float64
	Units   string
	out     chan<- Value
}

func NewFlowMeter(pin *Pin, opts ...func(InputDevice) error) (InputDevice, error) {
	pin.Direction = "in"
	pin.Value = true
	gpio, err := newGPIO(pin)
	if err != nil {
		return nil, err
	}
	minSpan := getMinSpan(pin)

	return &FlowMeter{
		GPIO:    gpio,
		Value:   0.0,
		Units:   pin.Units,
		MinSpan: minSpan,
	}, nil
}

func (f *FlowMeter) wait(err chan<- error) {
	for {
		e := f.GPIO.Wait()
		if e != nil {
			err <- e
			continue
		}
		t1 := f.ts
		f.ts = time.Now()
		if t1.Year() == 1 {
			continue
		}
		span := f.ts.Sub(t1).Seconds()
		if span < f.MinSpan {
			continue
		}
		f.value = f.Value / float64(span)
		f.out <- Value{
			Value: f.value,
			Units: f.Units,
		}
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

func getMinSpan(pin *Pin) float64 {
	v, ok := pin.Args["min_span"].(float64)
	if !ok {
		return 0.1
	}
	return v
}
