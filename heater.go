package gogadgets

import (
	"time"
	"fmt"
)

//Heater represnts an electic heating element.  It
//provides a way to heat up something to a target
//temperature. In order to use this there must be
//a thermometer in the same Location.
type Heater struct {
	target   float64
	currentTemp float64
	duration time.Duration
	status   bool
	pwm      OutputDevice
	doPWM    bool
	update   chan Message
	start     chan bool
	watching bool
}

func NewHeater(pin *Pin) (OutputDevice, error) {
     	var h *Heater
	var err error
	p, err := NewPWM(pin)
	doPWM := pin.Args["pwm"] == "true"
	if err == nil {
		h = &Heater{
			pwm:    p,
			target:  100.0,
			doPWM:   doPWM,
			update: make(chan Message),
			start: make(chan bool),
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
		} else {
			h.target = 100.0
		}
	}
	h.status = true
	if !h.watching {
		go h.watchTemperature(h.update, h.start)
	}
	h.start <- true
	return nil
}

func (h *Heater) Status() interface{} {
	return h.status
}

func (h *Heater) Off() error {
	h.target = 0.0
	h.start <- false
	return nil
}

func (h *Heater) watchTemperature(update <-chan Message, start <-chan bool) {
	fmt.Println("watching")
	h.watching = true
	for {
		select {
		case msg := <-update:
			fmt.Println("msg", msg)
			h.readTemperature(msg)
		case s := <-start:
			h.status = s
			fmt.Println("got start message", s)
			if s {
				h.toggle()
			} else {
				h.pwm.Off()
			}
		}
	}
}

func (h *Heater) readTemperature(msg Message) {
	temp, ok := msg.Value.ToFloat()
	fmt.Println(temp, ok)
	if ok {
		h.currentTemp = temp
		fmt.Println("if status", h.status)
		if h.status {
			h.toggle()
		}
	}
}
	
func (h *Heater) toggle() {
	fmt.Println("toggle", h.pwm.Status(), h.status)
	if h.doPWM {
		duty := h.getDuty()
		val := &Value{Value:duty, Units:"%s"}
		h.pwm.On(val)
	} else {
		diff := h.target - h.currentTemp
		fmt.Println("diff", diff, h.target, h.currentTemp)
		if diff > 0 {
			h.pwm.On(nil)
		} else {
			h.pwm.Off()
		}
	}
}

//Once the heater approaches the target temperature the electricity
//is applied PWM style so the target temperature isn't overshot.
func (h *Heater) getDuty() int {
	diff := h.target - h.currentTemp
	duty := 100
	if diff >= 1.0 {
		duty = 50
	} else if diff > 0.0 {
		duty = 25
	} else if diff <= 0.0 {
		duty = 0
	}
	return duty
}
