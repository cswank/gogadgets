package devices

import (
	"time"
	"testing"
	"bitbucket.com/cswank/gogadgets/models"
)

type FakePoller struct {
	Poller
	val bool
}

func (f *FakePoller) Wait() (bool, error) {
	time.Sleep(100 * time.Millisecond)
	f.val = !f.val
	return f.val, nil
}

func TestSwitch(t *testing.T) {
	poller := &FakePoller{}
	s := &Switch{
		GPIO: poller,
		Value: 5.0,
		Units: "liters",
	}
	stop := make(chan bool)
	in := make(chan models.Value)
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
