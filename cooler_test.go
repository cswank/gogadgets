package gogadgets

import (
	"testing"
)


func TestCreateCooler(t *testing.T) {
	_ = Cooler{
		gpio: &FakeOutput{},
		target: 0.0,
	}
}
