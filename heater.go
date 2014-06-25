package gogadgets

import (
	"time"
)

//Heater represents an electic heating element.  It
//provides a way to heat up something to a target
//temperature. In order to use this there must be
//a thermometer in the same Location.
type Heater struct {
	onTime      time.Duration
	offTime     time.Duration
	toggleTime  time.Duration
	waitTime    time.Duration
	t1          time.Time
	target      float64
	currentTemp float64
	duration    time.Duration
	status      bool
	gpioStatus  bool
	doPWM       bool
	gpio        OutputDevice
	io          chan *Value
	update      chan *Message
	started     bool
}

func NewHeater(pin *Pin) (OutputDevice, error) {
	var h *Heater
	var err error
	var dev OutputDevice
	doPWM := pin.Args["pwm"] == true
	if pin.Frequency == 0 {
		pin.Frequency = 1
	}
	dev, err = NewGPIO(pin)
	if err == nil {
		h = &Heater{
			toggleTime: 100 * time.Hour,
			gpio:       dev,
			target:     100.0,
			doPWM:      doPWM,
			io:         make(chan *Value),
			update:     make(chan *Message),
		}
	}
	return h, err
}

func (h *Heater) Config() ConfigHelper {
	return ConfigHelper{
		PinType: "pwm",
		Units:   []string{"C", "F"},
		Pins:    Pins["pwm"],
	}
}

func (h *Heater) Update(msg *Message) {
	if h.status && msg.Name == "temperature" {
		h.update <- msg
	} else {
		h.readTemperature(msg)
	}
}

func (h *Heater) On(val *Value) error {
	h.status = true
	if !h.started {
		h.started = true
		go h.toggle(h.io, h.update)
	}
	if val == nil {
		val = &Value{Value: true}
	}
	h.io <- val
	return nil
}

func (h *Heater) Status() interface{} {
	return h.status
}

func (h *Heater) Off() error {
	if h.started {
		h.target = 0.0
		h.status = false
		h.io <- &Value{Value: false}
	}
	return nil
}

/*
The pwm drivers on beaglebone black seem to be
broken.  This function brings the same functionality
using gpio.
*/
func (h *Heater) toggle(value chan *Value, update chan *Message) {
	for {
		select {
		case val := <-value:
			switch v := val.Value.(type) {
			case float64:
				h.waitTime = 100 * time.Millisecond
				h.getTarget(val)
				h.setDuty()
				h.status = true
				h.gpioStatus = true
				h.gpio.On(nil)
				h.t1 = time.Now()
			case bool:
				h.waitTime = 100 * time.Hour
				if v == true {
					h.status = true
					h.gpio.On(nil)
				} else {

					h.gpio.Off()
					h.target = 1000.0
					h.status = false
				}
			}
		case m := <-update:
			h.readTemperature(m)
		case _ = <-time.After(h.waitTime):
			n := time.Now()
			diff := n.Sub(h.t1)
			if h.doPWM && diff > h.toggleTime {
				h.t1 = n
				if h.gpioStatus && h.offTime > 0.0 {
					h.toggleTime = h.offTime
					h.gpio.Off()
					h.gpioStatus = false
				} else if !h.gpioStatus && h.onTime > 0.0 {
					h.toggleTime = h.onTime
					h.gpio.On(nil)
					h.gpioStatus = true
				} else {
					h.toggleTime = h.offTime
					h.gpio.Off()
					h.gpioStatus = false
				}
			}
		}
	}
}

func (h *Heater) getTarget(val *Value) {
	if val != nil {
		t, ok := val.ToFloat()
		if ok {
			h.target = t
		}
	}
}

func (h *Heater) readTemperature(msg *Message) {
	temp, ok := msg.Value.ToFloat()
	if ok {
		h.currentTemp = temp
		if h.status {
			h.setDuty()
		}
	}
}

//Once the heater approaches the target temperature the electricity
//is applied PWM style so the target temperature isn't overshot.
//This functionality is geared towards heating up a tank of water
//and can be disabled if you are using this component to heat something
//else, like a house.
func (h *Heater) setDuty() {
	diff := h.target - h.currentTemp
	if diff <= 0.0 {
		h.onTime = 0
		h.offTime = 1 * time.Second
	} else if diff <= 1.0 {
		h.onTime = 1 * time.Second
		h.offTime = 3 * time.Second
	} else if diff <= 2.0 {
		h.onTime = 2 * time.Second
		h.offTime = 2 * time.Second
	} else {
		h.onTime = 4 * time.Second
		h.offTime = 0 * time.Second
	}
	if h.gpioStatus {
		h.toggleTime = h.onTime
	} else {
		h.toggleTime = h.offTime
	}
}
