package gogadgets

import (

)


type PWM struct {

}

func NewPWM(pin *Pin) (OutputDevice, error) {
	pwm := &PWM{}
	return pwm, nil
}

func (p *PWM) Update(msg *Message) {
	
}

func (p *PWM) On(val *Value) error {
	return nil
}

func (p *PWM) Status() interface{} {
	return false
}

func (p *PWM) Off() error {
	return nil
}
