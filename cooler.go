package gogadgets

import "time"

type Cooler struct {
	target     float64
	status     bool
	gpio       OutputDevice
	lastChange *time.Time
}

func NewCooler(pin *Pin) (OutputDevice, error) {
	var c *Cooler
	var err error
	g, err := newGPIO(pin)
	if err == nil {
		c = &Cooler{
			gpio:   g,
			target: 0.0,
		}
	}
	return c, err
}

func (c *Cooler) Commands(location, name string) *Commands {
	return nil
}

func (c *Cooler) Update(msg *Message) bool {
	now := time.Now()
	if c.lastChange != nil && now.Sub(*c.lastChange) < 120*time.Second {
		return false
	}

	var changed bool
	temperature, ok := msg.Value.Value.(float64)
	if ok && c.status {
		changed = true
		if temperature <= c.target {
			c.gpio.Off()
			c.lastChange = &now
		} else if temperature > c.target {
			c.gpio.On(nil)
			c.lastChange = &now
		}
	}
	return changed
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

func (c *Cooler) Status() map[string]bool {
	return c.gpio.Status()
}
