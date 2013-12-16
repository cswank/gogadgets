package devices

import (
	"testing"
	"time"
	"bitbucket.com/cswank/gogadgets/utils"
)

func TestGPIO(t *testing.T) {
	if !utils.FileExists("/sys/class/gpio/export") {
		return //not a beaglebone
	}
	g, err := NewGPIO(&Pin{Port:"9", Pin:"15", Direction:"out"})
	if err != nil {
		t.Error(err)
	}
	g.On(nil)
	time.Sleep(1 * time.Second)
	g.Off()
}
