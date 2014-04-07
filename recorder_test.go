package gogadgets

import (
	"testing"
	"time"
)

func TestCreate(t *testing.T) {
	pin := &Pin{
		Args: map[string]string{
			"host": "localhost",
			"db":   "brewery",
			"summarize": "10",
		},
	}
	r, err := NewRecorder(pin)
	if err != nil {
		t.Error(err, r)
	}
}


func TestSummarize(t *testing.T) {
	r := Recorder{
		duration: time.Hour,
		history: map[string]summary{},
	}
	for i := 1; i < 11; i ++ {
		msg := &Message{
			Sender: "me",
			Value: Value{Value:float64(i)},
		}
		r.summarize(msg)
	}
	if r.history["me"].v != 55.0 {
		t.Error(r.history["me"].v)
	}
}

func TestGetFilter(t *testing.T) {
	s := "lab led, lab temperature"
	f := getFilter(s)
	if len(f) != 2 {
		t.Error(f)
	}
	if f[0] != "lab led" {
		t.Error(f)
	}
	s = ""
	f = getFilter(s)
	if len(f) != 0 {
		t.Error(f)
	}
}

