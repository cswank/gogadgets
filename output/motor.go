package output

import (
	"bitbucket.org/cswank/gogadgets/models"
)

type Motor struct {
	gpioA  OutputDevice
	gpioB  OutputDevice
	pwm    OutputDevice
	status bool
}

func NewMotor(pin *models.Pin) (OutputDevice, error) {
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
	return &Motor{
		gpioA: gpioA,
		gpioB: gpioB,
		pwm:   pwm,
	}, nil
}

func (m *Motor) Config() models.ConfigHelper {
	// g := GPIO{}
	// gpio := g.Config()
	// p := PWM{}
	// pwm := p.Config()
	return models.ConfigHelper{}
}

func (m *Motor) Update(msg *models.Message) {

}

func (m *Motor) On(val *models.Value) error {
	if val == nil {
		val = &models.Value{Value: 100.0, Units: "%"}
	}
	v, ok := val.Value.(float64)
	if !ok {
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
