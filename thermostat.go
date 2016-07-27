package gogadgets

import (
	"fmt"
	"strings"
	"time"
)

type cmp func(float64, float64) bool

/*Thermostat is used for controlling a furnace.
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

	//minimum time between state changes
	timeout time.Duration

	status     bool
	gpios      map[string]OutputDevice
	lastChange *time.Time
	lastCmd    string
	cmp        map[string]cmp

	//the location + name id of the temperature sensor (must be in the same location)
	sensor string
}

func NewThermostat(pin *Pin) (OutputDevice, error) {
	p, ok := pin.Pins["heat"]
	if !ok {
		return nil, fmt.Errorf("invalid heat pin: %v", pin)
	}
	h, err := NewGPIO(&p)
	if err != nil {
		lg.Fatal(err)
	}

	p, ok = pin.Pins["cool"]
	if !ok {
		return nil, fmt.Errorf("invalid cool pin: %v", pin)
	}

	c, err := NewGPIO(&p)
	if err != nil {
		lg.Fatal(err)
	}

	p, ok = pin.Pins["fan"]
	if !ok {
		return nil, fmt.Errorf("invalid fan pin: %v", pin)
	}

	f, err := NewGPIO(&p)
	if err != nil {
		lg.Fatal(err)
	}

	return &Thermostat{
		gpios: map[string]OutputDevice{
			"heat": h,
			"cool": c,
			"fan":  f,
		},
		cmp: map[string]cmp{
			"heat": func(x, y float64) bool { return x >= y },
			"cool": func(x, y float64) bool { return x < y },
		},
		sensor:  pin.Args["sensor"].(string),
		timeout: getTimeout(pin.Args),
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

func getTimeout(args map[string]interface{}) time.Duration {
	to := 5 * time.Minute

	i, ok := args["timeout"]
	if !ok {
		return to
	}

	s, ok := i.(string)
	if !ok {
		return to
	}

	if x, err := time.ParseDuration(s); err == nil {
		to = x
	}
	return to
}

func (t *Thermostat) Config() ConfigHelper {
	return ConfigHelper{
		PinType: "gpio",
		Units:   []string{"C", "F"},
		Pins:    Pins["gpio"],
	}
}

func (t *Thermostat) Update(msg *Message) bool {
	if msg.Sender != t.sensor {
		return false
	}
	now := time.Now()

	if t.lastChange != nil && now.Sub(*t.lastChange) < t.timeout {
		return false
	}

	var changed bool
	temperature, ok := msg.Value.Value.(float64)
	if t.status && ok && (t.lastCmd == "heat" || t.lastCmd == "cool") {
		changed = true
		gpio := t.gpios[t.lastCmd]
		fan := t.gpios["fan"]
		t.lastChange = &now
		if t.cmp[t.lastCmd](temperature, t.target) {
			gpio.Off()
			fan.Off()
		} else {
			gpio.On(nil)
			fan.On(nil)
		}
	}
	return changed
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
	t.lastChange = nil
	t.target = tar
	t.status = true
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
