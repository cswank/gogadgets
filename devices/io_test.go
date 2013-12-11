package devices

import (
	"bitbucket.com/cswank/gogadgets/pins"
	"testing"
)

func TestNewOutputDevice(t *testing.T) {
	pin := &.pins.Pin{Type: "gpio", Port: "9", Pin: "15"}
	d, err := NewOutputDevice(pin)
	if err != nil {
		t.Error(err)
	}
}
