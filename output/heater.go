package output

import (
	"bitbucket.org/cswank/gogadgets/models"
	"time"
)

//Heater represnts an electic heating element.  It
//provides a way to heat up something to a target
//temperature. In order to use this there must be
//a thermometer in the same Location.
type Heater struct {
	target      float64
	currentTemp float64
	duration    time.Duration
	status      bool
	doPWM       bool
	pwm         OutputDevice
	update      chan models.Message
	start       chan bool
	watching    bool
}

func NewHeater(pin *models.Pin) (OutputDevice, error) {
	var h *Heater
	var err error
	var d OutputDevice
	doPWM := pin.Args["pwm"] == true
	d, err = NewPWM(pin)
	if err == nil {
		h = &Heater{
			pwm:    d,
			target: 100.0,
			doPWM:  doPWM,
			update: make(chan models.Message),
			start:  make(chan bool),
		}
	}
	return h, err
}

func (h *Heater) Config() models.ConfigHelper {
	return models.ConfigHelper{
		Fields: map[string][]string{
			"port": []string{},
			"pin":  []string{},
		},
		Units: []string{"C", "F"},
	}
}

func (h *Heater) Update(msg *models.Message) {
	if h.status && msg.Name == "temperature" {
		h.update <- *msg
	}
}

func (h *Heater) On(val *models.Value) error {
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
		h.watching = true
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

func (h *Heater) watchTemperature(update <-chan models.Message, start <-chan bool) {
	for {
		select {
		case msg := <-update:
			h.readTemperature(msg)
		case s := <-start:
			h.status = s
			if s {
				h.toggle()
			} else {
				h.pwm.Off()
			}
		}
	}
}

func (h *Heater) readTemperature(msg models.Message) {
	temp, ok := msg.Value.ToFloat()
	if ok {
		h.currentTemp = temp
		if h.status {
			h.toggle()
		}
	}
}

func (h *Heater) toggle() {
	if h.doPWM {
		duty := h.getDuty()
		val := &models.Value{Value: duty, Units: "%"}
		h.pwm.On(val)
	} else {
		diff := h.target - h.currentTemp
		if diff > 0 {
			h.pwm.On(nil)
		} else {
			h.pwm.Off()
		}
	}
}

//Once the heater approaches the target temperature the electricity
//is applied PWM style so the target temperature isn't overshot.
func (h *Heater) getDuty() float64 {
	diff := h.target - h.currentTemp
	duty := 100.0
	if diff <= 0.0 {
		duty = 0.0
	} else if diff <= 1.0 {
		duty = 25.0
	} else if diff <= 2.0 {
		duty = 50.0
	}
	return duty
}
