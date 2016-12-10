package gogadgets

import (
	"fmt"

	"go.bug.st/serial.v1"
)

//RangeFinder represents a MaxSonar-EZ device.
type RangeFinder struct {
	//used to ranging on/off
	gpio *GPIO

	//for reading the range data
	port serial.Port

	value bool
	units string
}

func NewRangeFinder(pin *Pin) (InputDevice, error) {
	p, ok := pin.Args["port"].(string)
	if !ok {
		return nil, fmt.Errorf(`unable to create serial port for RangeFinder, pin.Args["port"] should be the path to a serial device`)
	}
}

func (r *RangeFinder) Config() ConfigHelper {
	return ConfigHelper{}
}

func (r *RangeFinder) GetValue() *Value {
	return &Value{
		Value: r.value,
		Units: r.units,
	}
}

func (r *RangeFinder) on() {
	r.gpio.On(nil)
}

func (r *RangeFinder) off() {
	r.gpio.Off()
}

func (r *RangeFinder) Start(in <-chan Message, out chan<- Value) {

}
