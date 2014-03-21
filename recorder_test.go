package gogadgets

import (
	"testing"
)

func TestCreate(t *testing.T) {
	pin := &Pin{
		Args: map[string]string{
			"host": "localhost",
			"db": "brewery",
		},
	}
	r, err := NewRecorder(pin)
	if err != nil {
		t.Error(err, r)
	}
}
