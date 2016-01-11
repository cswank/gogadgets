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
	sensor     string //the location + name id of the temperature sensor (must be in the same location)
}

func NewThermostat(pin *Pin) (OutputDevice, error) {
	var t *Thermostat
	var err error
	g, err := NewGPIO(pin)
	var c cmp

	var h, l float64
	if pin.Args["type"] == "cooler" {
		l = pin.Args["high"].(float64)
		h = pin.Args["low"].(float64)
		c = func(x, y float64) bool {
			return x <= y
		}
	} else {
		h = pin.Args["high"].(float64)
		l = pin.Args["low"].(float64)
		c = func(x, y float64) bool {
			return x >= y
		}
	}

	t = &Thermostat{
		gpio:       g,
		highTarget: h,
		lowTarget:  l,
		cmp:        c,
		sensor:     pin.Args["sensor"].(string),
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
	if msg.Sender != t.sensor {
		return
	}
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
