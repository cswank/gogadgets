package io

import (
	"testing"
)

func TestGPOutput(t *testing.T) {
	g, err := NewGPOutput("9", "15")
	if err != nil {
		t.Error(err)
	}
	g.On()
	time.Sleep(1 * time.Second)
	g.Off()
}
