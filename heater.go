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
	duration time.Duration
	status   bool
	pwm      OutputDevice
	doPWM    bool
	update   chan Message
	stop     chan bool
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
			stop: make(chan bool),
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
	fmt.Println("heater on", val)
	if val != nil {
		target, ok := val.ToFloat()
		if ok {
			h.target = target
		}
	}
	h.status = true
	if h.target > 0.0 {
		h.pwm.On(nil)
	}		
	if !h.watching {
		go h.watchTemperature(h.update, h.stop)
	}
	return nil
}

func (h *Heater) Status() interface{} {
	return h.status
}

func (h *Heater) Off() error {
	fmt.Println("heater is getting turned off")
	if h.status {
		h.stop <- true
	}
	h.target = 0.0
	return nil
}

func (h *Heater) watchTemperature(update <-chan Message, stop <-chan bool) {
	fmt.Println("watching")
	h.watching = true
	for {
		select {
		case msg := <-update:
			fmt.Println("msg", msg)
			h.readTemperature(msg)
		case <-stop:
			fmt.Println("got stop message")
			h.pwm.Off()
			h.status = false
		}
	}
}

func (h *Heater) readTemperature(msg Message) {
	temp, ok := msg.Value.ToFloat()
	fmt.Println(temp, ok)
	if ok {
		if h.status {
			h.toggle(temp)
		}
	}
}
	
func (h *Heater) toggle(temp float64) {
	fmt.Println("toggle", h.pwm.Status())
	if h.doPWM {
		duty := h.getDuty(temp)
		val := &Value{Value:duty, Units:"%s"}
		h.pwm.On(val)
	} else {
		diff := h.target - temp
		fmt.Println(diff, h.target, temp)
		if diff > 0 {
			h.pwm.On(nil)
		} else {
			h.pwm.Off()
		}
	}
}

//Once the heater approaches the target temperature the electricity
//is applied PWM style so the target temperature isn't overshot.
func (h *Heater) getDuty(temp float64) int {
	diff := h.target - temp
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
