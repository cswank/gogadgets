package gogadgets

import (
	"fmt"
	"log"
	"sync"
	"time"
)

/*
Alarm (when 'on') turns on a gpio when
the a particular event happens.  It turns
off after a set amount of time, or when turned
off.
*/
type Alarm struct {
	gpio     OutputDevice
	events   map[string]bool
	duration time.Duration
	status   bool
	lock     sync.Mutex
}

func NewAlarm(pin *Pin) (OutputDevice, error) {
	duration := getAlarmDuration(pin.Args["duration"])
	events, err := getAlarmEvents(pin.Args["events"])
	if err != nil {
		return nil, err
	}

	pin.Direction = "out"
	gpio, err := NewGPIO(pin)
	if err != nil {
		return nil, err
	}

	return &Alarm{
		gpio:     gpio,
		duration: duration,
		events:   events,
	}, nil
}

func (a *Alarm) Commands(location, name string) *Commands {
	return nil
}

func (a *Alarm) Config() ConfigHelper {
	return ConfigHelper{}
}

func (a *Alarm) Update(msg *Message) bool {
	if !a.status {
		return false
	}

	state, ok := a.events[msg.Sender]
	if !ok {
		return false
	}

	if msg.Value.Value.(bool) == state {
		a.lock.Lock()
		a.gpio.On(nil)
		a.lock.Unlock()
		go func() {
			time.Sleep(a.duration)
			a.lock.Lock()
			a.gpio.Off()
			a.lock.Unlock()
		}()
		return true
	}

	return false
}

func (a *Alarm) On(val *Value) error {
	a.status = true
	return nil
}

func (a *Alarm) Status() map[string]bool {
	return map[string]bool{
		"gpio": a.gpio.Status()["gpio"],
	}
}

func (a *Alarm) Off() error {
	a.status = false
	a.lock.Lock()
	a.gpio.Off()
	a.lock.Unlock()
	return nil
}

func getAlarmEvents(e interface{}) (map[string]bool, error) {
	tmp, ok := e.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("could not parse alarm events %v (should be map[string]bool)\n", e)
	}

	m := map[string]bool{}
	for k, v := range tmp {
		b, ok := v.(bool)
		if !ok {
			return nil, fmt.Errorf("could not parse alarm events %v (should be map[string]bool)\n", e)
		}
		m[k] = b
	}
	return m, nil
}

func getAlarmDuration(a interface{}) time.Duration {
	s, ok := a.(string)
	if !ok {
		s = "5m"
	}

	d, err := time.ParseDuration(s)
	if err != nil {
		log.Printf("WARNING: could not parse alarm duration %v, defaulting to 5m, err: %s\n", a, err)
		d, _ = time.ParseDuration("5m")
	}
	return d
}
