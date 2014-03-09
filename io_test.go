package gogadgets

import (
	"bitbucket.org/cswank/gogadgets/utils"
	"testing"
)

func TestNewOutputDevice(t *testing.T) {
	if !utils.FileExists("/sys/class/gpio/export") {
		return //not a beaglebone
	}
	pin := &Pin{Type: "gpio", Port: "9", Pin: "15"}
	d, err := NewOutputDevice(pin)
	if err != nil {
		t.Error(err, d)
	}
}
