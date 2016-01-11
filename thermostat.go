package gogadgets

import "time"

type cmp func(float64, float64) bool

type Thermostat struct {
	target     float64
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
			gpio:   g,
			target: 0.0,
			cmp:    c,
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
	if ok && t.status {
		if t.cmp(temperature, t.target) {
			t.gpio.Off()
			t.lastChange = &now
		} else {
			t.gpio.On(nil)
			t.lastChange = &now
		}
	}
}

func (t *Thermostat) On(val *Value) error {
	if val != nil {
		target, ok := val.Value.(float64)
		if ok {
			t.target = target
		}
	}
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
