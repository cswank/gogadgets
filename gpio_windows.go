package gogadgets

type GPIO struct{}

func NewGPIO(pin *Pin) (OutputDevice, error) {
	return nil, nil
}

func (g *GPIO) Init() error {
	return nil
}

func (g *GPIO) Update(msg *Message) {}

func (g *GPIO) On(val *Value) error {
	return nil
}

func (g *GPIO) Status() map[string]bool {
	return nil
}

func (g *GPIO) Off() error {
	return nil
}

func (g *GPIO) writeValue(path, value string) error {
	return nil
}

func (g *GPIO) Wait() (bool, error) {
	return false, nil
}
