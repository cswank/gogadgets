package gogadgets

import (
	"time"
)

type Heater struct {
	OutputDevice
	target   float64
	current  float64
	duration time.Duration
	status   bool
	gpio     OutputDevice
	update   chan Message
	stop     chan bool
}

func NewHeater(pin *Pin) (OutputDevice, error) {
	var h *Heater
	var err error
	g, err := NewGPIO(pin)
	if err == nil {
		h = &Heater{
			gpio:    g,
			current: 0.0,
			target:  100.0,
		}
	}
	return h, err
}

func (h *Heater) Update(msg *Message) {
	if h.status && msg.Name == "temperature" {
		h.update <- *msg
	}
}

func (h *Heater) On(val *Value) error {
	if val != nil {
		target, ok := val.ToFloat()
		if ok {
			h.target = target
		}
	}
	h.status = true
	h.update = make(chan Message)
	h.stop = make(chan bool)
	go h.watchTemperature(h.update, h.stop)
	return nil
}

func (h *Heater) Status() interface{} {
	return h.status
}

func (h *Heater) Off() error {
	if h.status {
		h.stop <- true
	}
	return nil
}

func (h *Heater) watchTemperature(update <-chan Message, stop <-chan bool) {
	h.gpio.On(nil)
	keepGoing := true
	h.duration = time.Duration(60 * time.Second)
	for keepGoing {
		select {
		case msg := <-update:
			current, ok := msg.Value.ToFloat()
			if ok {
				h.current = current
				h.toggle()
			}
		case <-stop:
			h.gpio.Off()
			keepGoing = false
		case <-time.After(h.duration):
			h.toggle()
		}
	}
}

func (h *Heater) toggle() {
	on, off := h.getDurations()
	status := h.gpio.Status().(bool)
	if on == 0 && off != 0 {
		h.gpio.Off()
		h.duration = off
	} else if status && off != 0 {
		h.gpio.Off()
		h.duration = off
	} else if !status && on != 0 {
		h.gpio.On(nil)
		h.duration = on
	}
}

func (h *Heater) getDurations() (on time.Duration, off time.Duration) {
	diff := h.target - h.current
	if diff >= 2.0 {
		on = time.Duration(60.0 * float64(time.Second))
	} else if diff >= 1.0 {
		on = time.Duration(0.5 * float64(time.Second))
		off = time.Duration(0.5 * float64(time.Second))
	} else if diff > 0.0 {
		on = time.Duration(0.25 * float64(time.Second))
		off = time.Duration(0.75 * float64(time.Second))
	} else if diff <= 0.0 {
		off = time.Duration(60.0 * float64(time.Second))
	}
	return on, off
}
