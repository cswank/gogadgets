package gogadgets

import (
	"fmt"
)

type Door struct {
	motor OutputDevice
	poller InputDevice
	started bool
	lastDirection bool
}

func NewDoor(pin *Pin) (OutputDevice, error) {
	motor, err := NewMotor(pin)
	p = pin.Pins["poller"]
	poller, err := NewGPIO(&p)
	if err != nil {
		return nil, err
	}
	return &Door{
		motor: motor,
		poller: poller,
	}, nil
}

func (d *Door) On(val *Value) error {
	if !d.started {
		d.start()
	}
	d.motor.On(val)
	return nil
}

func (d *Door) Status() interface{} {
	return d.motor.Status()
}

func (d *Door) Off() error {
	return d.motor.Off()
}

func (d *Door) start() {
	d.started = true
}

func (d *Door) wait() {
	
}

func (d *Door) Update(msg *Message) {
	
}
