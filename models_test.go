package gogadgets

import (
	"testing"
)

func TestIsMine(t *testing.T) {
	g := NewOutputGadget("lab", "led")

	msg := &Message{
		Command: "turn on lab led",
	}
	if !g.isMine(msg) {
		t.Error("should have been mine")
	}
}

