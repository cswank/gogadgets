package gogadgets

import (
	"fmt"
	"log"
	"strings"
)

type cmp func(float64, float64, float64, bool) bool

/*
Thermostat is used for controlling a furnace.
Configure a thermostat like:

		{
		    "host": "http://192.168.1.18:6111",
		    "gadgets": [
		        {
		            "location": "home",
		            "name": "temperature",
		            "pin": {
		                "type": "thermometer",
		                "OneWireId": "28-0000041cb544",
		                "Units": "F"
		            }
		        },
		        {
		            "location": "home",
		            "name": "furnace",
		            "pin": {
		                "type": "thermostat",
	                    "pins": {
	                        "heat": {
	                            "platform": "rpi",
		                        "pin": "11",
	                            "direction": "out"
	                        },
	                        "cool": {
	                            "platform": "rpi",
		                        "pin": "13",
	                            "direction": "out"
	                        },
	                        "fan": {
	                            "platform": "rpi",
		                        "pin": "15",
	                            "direction": "out"
	                        }
	                    },
		                "args": {
		                    "sensor": "home temperature",
	                        "timeout": "5m"
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
	target float64

	hysteresis float64

	status          bool
	gpios           map[string]OutputDevice
	lastCmd         string
	lastTemperature *float64
	cmp             map[string]cmp

	//the location + name id of the temperature sensor (must be in the same location)
	sensor string
}

func NewThermostat(pin *Pin) (OutputDevice, error) {
	p, ok := pin.Pins["heat"]
	if !ok {
		return nil, fmt.Errorf("invalid heat pin: %v", pin)
	}
	h, err := newGPIO(&p)
	if err != nil {
		log.Fatal(err)
	}

	p, ok = pin.Pins["cool"]
	if !ok {
		return nil, fmt.Errorf("invalid cool pin: %v", pin)
	}

	c, err := newGPIO(&p)
	if err != nil {
		log.Fatal(err)
	}

	p, ok = pin.Pins["fan"]
	if !ok {
		return nil, fmt.Errorf("invalid fan pin: %v", pin)
	}

	f, err := newGPIO(&p)
	if err != nil {
		return nil, err
	}

	hy, err := getHysteresis(pin.Args)
	if err != nil {
		return nil, err
	}

	return &Thermostat{
		gpios: map[string]OutputDevice{
			"heat": h,
			"cool": c,
			"fan":  f,
		},
		cmp: map[string]cmp{
			"heat": func(actual, target, hysteresis float64, status bool) bool {
				if status {
					return actual >= (target + hysteresis)
				}
				return actual >= (target - hysteresis)
			},
			"cool": func(actual, target, hysteresis float64, status bool) bool {
				if status {
					return actual <= (target - hysteresis)
				}
				return actual <= (target + hysteresis)
			},
		},
		sensor:     pin.Args["sensor"].(string),
		hysteresis: hy,
	}, nil
}

func (t *Thermostat) Commands(location, name string) *Commands {
	return &Commands{
		On: []string{
			fmt.Sprintf("heat %s", location),
			fmt.Sprintf("cool %s", location),
		},
		Off: []string{
			fmt.Sprintf("turn off %s", name),
		},
	}
}

func (t *Thermostat) Update(msg *Message) bool {
	if msg.Sender != t.sensor {
		return false
	}

	var changed bool
	temperature, ok := msg.Value.Value.(float64)
	if t.status && ok && (t.lastCmd == "heat" || t.lastCmd == "cool") {
		t.lastTemperature = &temperature
		changed = t.checkTemperature()
	}
	return changed
}

func (t *Thermostat) checkTemperature() bool {
	if t.lastTemperature == nil {
		return false
	}

	gpio := t.gpios[t.lastCmd]
	st := gpio.Status()["gpio"]
	if t.cmp[t.lastCmd](*t.lastTemperature, t.target, t.hysteresis, st) {
		gpio.Off()
		return st
	} else {
		gpio.On(nil)
		return !st
	}
}

func (t *Thermostat) On(val *Value) error {
	if val == nil {
		return nil
	}
	tar, ok := val.Value.(float64)
	if !ok {
		return nil
	}
	parts := strings.Split(val.Cmd, " ")
	if len(parts) == 0 {
		return nil
	}
	t.lastCmd = parts[0]
	t.target = tar
	t.status = true
	t.checkTemperature()
	return nil
}

func (t *Thermostat) Off() error {
	if t.status {
		t.status = false
		t.gpios["heat"].Off()
		t.gpios["cool"].Off()
		t.gpios["fan"].Off()
	}
	return nil
}

func (t *Thermostat) Status() map[string]bool {
	m := map[string]bool{}
	for key, val := range t.gpios {
		m[key] = val.Status()["gpio"]
	}
	return m
}

func getHysteresis(args map[string]interface{}) (float64, error) {
	out := float64(5)

	i, ok := args["hysteresis"]
	if !ok {
		return out, nil
	}

	h, ok := i.(float64)
	if !ok {
		return out, fmt.Errorf("hysteresis is not float64")
	}

	return h, nil
}
