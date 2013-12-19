package utils

import (
	"testing"
)

func TestFileExists(t *testing.T) {
	if !FileExists("/etc/hosts") {
		t.Error("should exist")
	}

	if FileExists("/nowhere/stuff") {
		t.Error("should exist")
	}
}
