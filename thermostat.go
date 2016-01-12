package gogadgets

import "time"

type cmp func(float64, float64) bool

/*
Configure a thermostat like:

	{
	    "host": "http://192.168.1.30:6111",
	    "gadgets": [
	        {
	            "location": "the lab",
	            "name": "temperature",
	            "pin": {
	                "type": "thermometer",
	                "OneWireId": "28-0000041cb544",
	                "Units": "F"
	            }
	        },
	        {
	            "location": "the lab",
	            "name": "heater",
	            "pin": {
	                "type": "thermostat",
	                "port": "8",
	                "pin": "11",
	                "args": {
	                    "type": "heater",
	                    "sensor": "the lab temperature",
	                    "high": 150.0,
	                    "low": 120.0
	                }
	            }
	        }
	    ]
	}

With this config the thermostat will react to temperatures from
'the lab temperature' (which is the location + name of the thermometer)
and turn on the gpio if the temperature is > 120.0, turn and turn it
off when the temperature > 150.0.

If you set args.type = "cooler" then it will start cooling when the
temperature gets above 150, and stop cooling when the temperature gets
below 120.
*/
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
