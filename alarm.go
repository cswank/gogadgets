package gogadgets

import (
	"fmt"
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
	out      map[string]OutputDevice
	events   map[string]bool
	duration time.Duration
	delay    time.Duration
	status   bool
	lock     sync.Mutex
	ch       chan bool
}

func NewAlarm(pin *Pin) (OutputDevice, error) {
	duration, err := getAlarmDuration(pin.Args["duration"])
	if err != nil {
		return nil, err
	}

	delay, err := getAlarmDuration(pin.Args["delay"])
	events, err := getAlarmEvents(pin.Args["events"])
	if err != nil {
		return nil, err
	}

	a := &Alarm{
		duration: duration,
		delay:    delay,
		events:   events,
		ch:       make(chan bool),
	}

	go a.trigger()

	return a, nil
}

func (a *Alarm) WithOutput(out map[string]OutputDevice) {
	a.out = out
}

func (a *Alarm) Commands(location, name string) *Commands {
	return nil
}

func (a *Alarm) Config() ConfigHelper {
	return ConfigHelper{}
}

func (a *Alarm) trigger() {
	delay := time.After(100000 * time.Hour)
	duration := time.After(100000 * time.Hour)

	for {
		select {
		case <-duration:
			a.lock.Lock()
			for _, out := range a.out {
				out.Off()
			}
			a.lock.Unlock()
			duration = time.After(100000 * time.Hour)
		case <-delay:
			a.lock.Lock()
			for _, out := range a.out {
				out.On(nil)
			}
			a.lock.Unlock()
			duration = time.After(a.duration)
			delay = time.After(100000 * time.Hour)
		case b := <-a.ch:
			if !b {
				delay = time.After(100000 * time.Hour)
				duration = time.After(100000 * time.Hour)
			} else {
				delay = time.After(a.delay)
			}
		}
	}
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
		a.ch <- true
	} else {
		a.ch <- false
	}

	return false
}

func (a *Alarm) On(val *Value) error {
	a.status = true
	return nil
}

func (a *Alarm) Status() map[string]bool {
	m := map[string]bool{}
	for name, out := range a.out {
		m[name] = out.Status()[name]
	}

	return m
}

func (a *Alarm) Off() error {
	a.status = false
	a.ch <- false
	a.lock.Lock()
	for _, out := range a.out {
		out.Off()
	}
	a.lock.Unlock()
	return nil
}

func getAlarmEvents(e interface{}) (map[string]bool, error) {
	m, ok := e.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("could not parse alarm events %v (should be map[string]bool)\n", e)
	}

	out := map[string]bool{}
	for k, v := range m {
		b, ok := v.(bool)
		if !ok {
			return nil, fmt.Errorf("could not parse alarm events %v (should be map[string]bool)\n", e)
		}
		out[k] = b
	}
	return out, nil
}

func getAlarmDuration(a interface{}) (time.Duration, error) {
	var d time.Duration
	s, ok := a.(string)
	if !ok {
		return d, fmt.Errorf("Could not parse alarm duration %v", a)
	}

	return time.ParseDuration(s)
}
