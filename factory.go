package gogadgets

import (
	"errors"
)

type AppFactory struct {
	inputFactories  map[string]InputDeviceFactory
	outputFactories map[string]OutputDeviceFactory
}

func NewAppFactory() *AppFactory {
	a := &AppFactory{
		inputFactories: map[string]InputDeviceFactory{
			"thermometer": NewThermometer,
			"switch":      NewSwitch,
		},
		outputFactories: map[string]OutputDeviceFactory{
			"gpio":     NewGPIO,
			"heater":   NewHeater,
			"cooler":   NewCooler,
			"recorder": NewRecorder,
		},
	}
	return a
}

//Each input and output device has a config method that returns a Pin with
//the required fields poplulated with helpful values.
func GetTypes() map[string]ConfigHelper {
	t := Thermometer{}
	s := Switch{}
	g := GPIO{}
	h := Heater{}
	c := Cooler{}
	r := Recorder{}
	return map[string]ConfigHelper{
		"thermometer": t.Config(),
		"switch":      s.Config(),
		"gpio":        g.Config(),
		"heater":      h.Config(),
		"cooler":      c.Config(),
		"recorder":    r.Config(),
	}
}

func (f *AppFactory) RegisterInputFactory(name string, factory InputDeviceFactory) {
	f.inputFactories[name] = factory
}

func (f *AppFactory) RegisterOutputFactory(name string, factory OutputDeviceFactory) {
	f.outputFactories[name] = factory
}

func (f *AppFactory) GetApp() (a *App, err error) {
	return a, err
}

type InputDeviceFactory func(pin *Pin) (InputDevice, error)

//Inputdevices are started as goroutines by the Gadget
//that contains it.
type InputDevice interface {
	Start(<-chan Message, chan<- Value)
	GetValue() *Value
	Config() ConfigHelper
}

type Poller interface {
	Wait() (bool, error)
}

func NewInputDevice(pin *Pin) (dev InputDevice, err error) {
	if pin.Type == "thermometer" {
		dev, err = NewThermometer(pin)
	} else if pin.Type == "switch" {
		dev, err = NewSwitch(pin)
	} else {
		err = errors.New("invalid pin type")
	}
	return dev, err
}

type OutputDeviceFactory func(pin *Pin) (OutputDevice, error)

//Outputdevices turn things on and off.  Currently the
//only
type OutputDevice interface {
	On(val *Value) error
	Off() error
	Update(msg *Message)
	Status() interface{}
	Config() ConfigHelper
}

func NewOutputDevice(pin *Pin) (dev OutputDevice, err error) {
	if pin.Type == "gpio" {
		dev, err = NewGPIO(pin)
	} else if pin.Type == "heater" {
		dev, err = NewHeater(pin)
	} else if pin.Type == "cooler" {
		dev, err = NewCooler(pin)
	} else if pin.Type == "recorder" {
		dev, err = NewRecorder(pin)
	} else if pin.Type == "pwm" {
		dev, err = NewPWM(pin)
	} else if pin.Type == "motor" {
		dev, err = NewMotor(pin)
	} else {
		dev, err = nil, errors.New("invalid pin type")
	}
	return dev, err
}
