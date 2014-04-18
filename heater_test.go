package gogadgets

import (
	"os"
	"testing"
)

func init() {
	os.Setenv("PWMTESTDEVICEMODE", "perm")
}

func getMessage(val float64) *Message {
	return &Message{
		Name: "temperature",
		Value: Value{
			Value: val,
		},
	}
}

func TestHeater(t *testing.T) {
	p := &Pin{
		Type: "heater",
		Port: "8",
		Pin: "13",
		Frequency: 1,
	}
	d, err := NewHeater(p)
	if err != nil {
		t.Error(err, d)
	}
}

