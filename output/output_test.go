package output

import (
	"bitbucket.org/cswank/gogadgets/utils"
	"bitbucket.org/cswank/gogadgets/models"
	"testing"
)

func TestNewOutputDevice(t *testing.T) {
	if !utils.FileExists("/sys/class/gpio/export") {
		return //not a beaglebone
	}
	pin := &models.Pin{Type: "gpio", Port: "9", Pin: "15"}
	d, err := NewOutputDevice(pin)
	if err != nil {
		t.Error(err, d)
	}
}
