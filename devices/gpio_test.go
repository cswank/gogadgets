package devices

import (
	"fmt"
	"testing"
	"time"
	"bitbucket.com/cswank/gogadgets/utils"
	"bitbucket.com/cswank/gogadgets/models"
)

func TestGPIO(t *testing.T) {
	if !utils.FileExists("/sys/class/gpio/export") {
		return //not a beaglebone
	}
	g, err := NewGPIO(&models.Pin{Port:"9", Pin:"15", Direction:"out"})
	if err != nil {
		t.Error(err)
	}
	g.On(nil)
	time.Sleep(1 * time.Second)
	g.Off()
}


func TestGPIOWait(t *testing.T) {
	g, err := NewGPIO(&models.Pin{Port:"9", Pin:"16", Direction:"in", Edge:"rising"})
	if err != nil {
		t.Error(err)
	}
	gIn, _ := NewGPIO(&models.Pin{Port:"9", Pin:"15", Direction:"out"})
	go func() {
		gIn.Off()
		time.Sleep(1 * time.Second)
		fmt.Println("turning on gpio")
		gIn.On(nil)
	}()
	fmt.Println("wait()")
	val, err := g.Wait()
	if err != nil {
		t.Error(err)
	}
	if val != true {
		t.Error("should have got true")
	}
	gIn.Off()
}

