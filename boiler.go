package gogadgets

import "time"

/*
Configure a boiler like:

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
	                "type": "boiler",
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

With this config the boiler will react to temperatures from
'the lab temperature' (which is the location + name of the thermometer)
and turn on the gpio if the temperature is > 120.0, turn and turn it
off when the temperature > 150.0.

If you set args.type = "cooler" then it will start cooling when the
temperature gets above 150, and stop cooling when the temperature gets
below 120.
*/
type Boiler struct {
	highTarget float64
	lowTarget  float64
	timeout    time.Duration
	status     bool
	gpio       OutputDevice
	lastChange *time.Time
	cmp        cmp
	sensor     string //the location + name id of the temperature sensor (must be in the same location)
}

func NewBoiler(pin *Pin) (OutputDevice, error) {
	var b *Boiler
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

	b = &Boiler{
		gpio:       g,
		highTarget: h,
		lowTarget:  l,
		cmp:        c,
		sensor:     pin.Args["sensor"].(string),
	}
	return b, err
}

func (b *Boiler) Commands(location, name string) *Commands {
	return nil
}

func (b *Boiler) Config() ConfigHelper {
	return ConfigHelper{
		PinType: "gpio",
		Units:   []string{"C", "F"},
		Pins:    Pins["gpio"],
	}
}

func (b *Boiler) Update(msg *Message) bool {
	if msg.Sender != b.sensor {
		return false
	}
	now := time.Now()
	// if b.lastChange != nil && now.Sub(*b.lastChange) < 120*time.Second {
	// 	return
	// }
	var ch bool
	temperature, ok := msg.Value.Value.(float64)
	if b.status && ok {
		ch = true
		if b.cmp(temperature, b.highTarget) {
			b.gpio.Off()
			b.lastChange = &now
		} else if b.cmp(b.lowTarget, temperature) {
			b.gpio.On(nil)
			b.lastChange = &now
		}
	}
	return ch
}

func (b *Boiler) On(val *Value) error {
	b.status = true
	b.gpio.On(nil)
	return nil
}

func (b *Boiler) Off() error {
	if b.status {
		b.status = false
		b.gpio.Off()
	}
	return nil
}

func (b *Boiler) Status() map[string]bool {
	return b.gpio.Status()
}
