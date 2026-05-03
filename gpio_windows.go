package gogadgets

import "os"

var (
	GPIO_DEV_PATH = "/sys/class/gpio"
	GPIO_DEV_MODE = os.ModeDevice
)

type GPIO struct{}

func NewGPIO(pin *Pin) (OutputDevice, error) {
	return nil, nil
}

func newGPIO(pin *Pin) (*GPIO, error) {
	return nil, nil
}

func (g *GPIO) Commands(location, name string) *Commands {
	return nil
}

func (g *GPIO) Update(msg *Message) bool {
	return false
}

func (g *GPIO) On(val *Value) error {
	return nil
}

func (g *GPIO) Status() map[string]bool {
	return nil
}

func (g *GPIO) Off() error {
	return nil
}

func (g *GPIO) Wait() error {
	return nil
}

func (g *GPIO) Close() error {
	return nil
}
