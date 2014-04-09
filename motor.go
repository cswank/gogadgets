package gogadgets

import (
	
)

type Motor struct {
	gpioA OutputDevice
	gpioB OutputDevice
	pwm   OutputDevice
	poller InputDevice
	status bool
	started bool
}

func NewMotor(pin *Pin) (OutputDevice, error) {
	p := pin.Pins["gpio_a"]
	gpioA, err := NewGPIO(&p)
	if err != nil {
		return nil, err
	}
	p = pin.Pins["gpio_b"]
	gpioB, err := NewGPIO(&p)
	if err != nil {
		return nil, err
	}
	p = pin.Pins["pwm"]
	pwm, err := NewPWM(&p)
	if err != nil {
		return nil, err
	}
	p = pin.Pins["poller"]
	poller, err := NewSwitch(&p)
	if err != nil {
		return nil, err
	}
	return &Motor{
		gpioA: gpioA,
		gpioB: gpioB,
		pwm: pwm,
		poller: poller,
	}, nil
}

func (m *Motor) Update(msg *Message) {
	
}

func (m *Motor) On(val *Value) error {
	v, ok := val.Value.(float64)
	if ! ok {
		return nil
	}
	if v < 0.0 {
		m.pwm.On(val)
		m.gpioA.Off()
		m.gpioB.On(nil)
	} else if v > 0.0 {
		m.pwm.On(val)
		m.gpioA.On(nil)
		m.gpioB.Off()
	} else {
		m.Off()
	}
	return nil
}

func (m *Motor) Status() interface{} {
	return m.status
}

func (m *Motor) Off() error {
	m.pwm.Off()
	m.gpioA.On(nil)
	m.gpioB.On(nil)
	return nil
}


func (m *Motor) start() {
	m.started = true
}
