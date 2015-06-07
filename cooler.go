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
	g, err := NewGPIO(pin)
	if err == nil {
		c = &Cooler{
			gpio:   g,
			target: 0.0,
		}
	}
	return c, err
}

func (c *Cooler) Config() ConfigHelper {
	return ConfigHelper{
		PinType: "gpio",
		Units:   []string{"C", "F"},
		Pins:    Pins["gpio"],
	}
}

func (c *Cooler) Update(msg *Message) {
	now := time.Now()
	if c.lastChange != nil && now.Sub(*c.lastChange) < 120*time.Second {
		return
	}
	temperature, ok := msg.Value.Value.(float64)
	if ok && c.status {
		if temperature <= c.target {
			c.gpio.Off()
			c.lastChange = &now
		} else if temperature > c.target {
			c.gpio.On(nil)
			c.lastChange = &now
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
