package gogadgets

import "time"

type cmp func(float64, float64) bool

type Thermostat struct {
	highTarget float64
	lowTarget  float64
	timeout    time.Duration
	status     bool
	gpio       OutputDevice
	lastChange *time.Time
	cmp        cmp
}

func NewThermostat(pin *Pin) (OutputDevice, error) {
	var t *Thermostat
	var err error
	g, err := NewGPIO(pin)
	var c cmp
	if pin.Args["type"] == "cooler" {
		c = func(x, y float64) bool {
			return x <= y
		}
	} else {
		c = func(x, y float64) bool {
			return x >= y
		}
	}
	if err == nil {
		t = &Thermostat{
			gpio:       g,
			highTarget: pin.Args["high"].(float64),
			lowTarget:  pin.Args["low"].(float64),
			cmp:        c,
		}
	}
	return t, err
}

func (t *Thermostat) Config() ConfigHelper {
	return ConfigHelper{
		PinType: "gpio",
		Units:   []string{"C", "F"},
		Pins:    Pins["gpio"],
	}
}

func (t *Thermostat) Update(msg *Message) {
	now := time.Now()
	// if t.lastChange != nil && now.Sub(*t.lastChange) < 120*time.Second {
	// 	return
	// }
	temperature, ok := msg.Value.Value.(float64)
	if t.status && ok {
		if t.cmp(temperature, t.highTarget) {
			t.gpio.Off()
			t.lastChange = &now
		} else if t.cmp(t.lowTarget, temperature) {
			t.gpio.On(nil)
			t.lastChange = &now
		}
	}
}

func (t *Thermostat) On(val *Value) error {
	t.status = true
	t.gpio.On(nil)
	return nil
}

func (t *Thermostat) Off() error {
	if t.status {
		t.status = false
		t.gpio.Off()
	}
	return nil
}

func (t *Thermostat) Status() interface{} {
	return t.status
}
