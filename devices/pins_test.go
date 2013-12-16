package devices

import (
	"testing"
)

func TestGPIO(t *testing.T) {
	if GPIO["8"]["7"] != "66" {
		t.Error(GPIO)
	}
	if GPIO["9"]["15"] != "48" {
		t.Error(GPIO)
	}
}
