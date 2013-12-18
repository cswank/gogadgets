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
		gpio: poller,
	}
	stop := make(chan bool)
	in := make(chan models.Value)
	go s.Start(stop, in)
	val := <-in
	if val.Value != true {
		t.Error("should have been true")
	}
	val = <-in
	if val.Value != false {
		t.Error("should have been false")
	}
	stop<- true
}
