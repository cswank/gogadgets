package gogadgets

import ()

type Cooler struct {
	target float64
	status bool
	gpio   OutputDevice
}

func NewCooler(pin *Pin) (OutputDevice, error) {
	var c *Cooler
	var err error
	g, err := NewGPIO(pin)
	if err == nil {
		c = &Cooler{
			gpio:   g,
			target: 0.0,
		}
	}
	return c, err
}

func (c *Cooler) Update(msg *Message) {
	temperature, ok := msg.Value.Value.(float64)
	if ok && c.status {
		if temperature <= c.target {
			c.gpio.Off()
		} else if temperature > c.target {
			c.gpio.On(nil)
		}
	}
}

func (c *Cooler) On(val *Value) error {
	if val != nil {
		target, ok := val.Value.(float64)
		if ok {
			c.target = target
		}
	}
	c.status = true
	c.gpio.On(nil)
	return nil
}

func (c *Cooler) Off() error {
	if c.status {
		c.status = false
		c.gpio.Off()
	}
	return nil
}

func (c *Cooler) Status() interface{} {
	return c.status
}
