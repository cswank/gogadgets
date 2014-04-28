package output

import (
	"testing"
)

type FakeOutput struct {
	OutputDevice
	on bool
}


func TestCreateCooler(t *testing.T) {
	_ = Cooler{
		gpio:   &FakeOutput{},
		target: 0.0,
	}
}
