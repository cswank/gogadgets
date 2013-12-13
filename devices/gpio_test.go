package devices

import (
	"testing"
	"time"
)

func TestGPIO(t *testing.T) {
	g, err := NewGPIO(&Pin{Port:"9", Pin:"15", Direction:"out"})
	if err != nil {
		t.Error(err)
	}
	g.On()
	time.Sleep(1 * time.Second)
	g.Off()
}
