package gogadgets

import (
	
)

type Motor struct {
	gpioA GPIO
	gpioB GPIO
	pwm   PWM
	input GPIO
	status bool
}

func NewMotor(pin *Pin) (OutputDevice, error) {
	var m *Motor
	var err error
	return m, err
}

func (m *Motor) Update(msg *Message) {
	
}

func (m *Motor) On(val *Value) error {
	return nil
}

func (m *Motor) Status() interface{} {
	return m.status
}

func (m *Motor) Off() error {
	return nil
}
