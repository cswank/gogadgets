package input

import (
	"time"
	"testing"
	"bitbucket.org/cswank/gogadgets/models"
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
		GPIO:      poller,
		Value:     5.0,
		TrueValue: 5.0,
		Units:     "liters",
	}
	out := make(chan models.Message)
	in := make(chan models.Value)
	go s.Start(out, in)
	val := <-in
	if val.Value.(float64) != 5.0 {
		t.Error("should have been 5.0", val)
	}
	val = <-in
	if val.Value.(float64) != 0.0 {
		t.Error("should have been 0.0", val)
	}
	out <- models.Message{
		Type: "command",
		Body: "shutdown",
	}
	v := s.GetValue()
	if v.Value.(float64) != 0.0 {
		t.Error("should have been 0.0", v)
	}
}

func TestBoolSwitch(t *testing.T) {
	poller := &FakePoller{}
	s := &Switch{
		GPIO:      poller,
		Value:     true,
		TrueValue: true,
	}
	out := make(chan models.Message)
	in := make(chan models.Value)
	go s.Start(out, in)
	val := <-in
	if val.Value != true {
		t.Error("should have been true", val)
	}
	val = <-in
	if val.Value != false {
		t.Error("should have been false", val)
	}
	out <- models.Message{
		Type: "command",
		Body: "shutdown",
	}
	v := s.GetValue()
	if v.Value != false {
		t.Error("should have been false", v)
	}
}
