package devices

import (
	"testing"
	"bitbucket.com/cswank/gogadgets/utils"
	"bitbucket.com/cswank/gogadgets/models"
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
