package gogadgets

import (
	"testing"
)

func TestSwitch(t *testing.T) {
	poller := &FakePoller{}
	s := &Switch{
		GPIO: poller,
		Value: 5.0,
		Units: "liters",
	}
	stop := make(chan bool)
	in := make(chan Value)
	go s.Start(stop, in)
	val := <-in
	if val.Value.(float64) != 5.0 {
		t.Error("should have been 5.0", val)
	}
	val = <-in
	if val.Value.(float64) != 0.0 {
		t.Error("should have been 0.0", val)
	}
	stop<- true
}
