package gogadgets

import (
	"errors"
)

type AppFactory struct {
	inputFactories  map[string]InputDeviceFactory
	outputFactories map[string]OutputDeviceFactory
}

var (
	inputFactories = map[string]InputDeviceFactory{
		"thermometer": NewThermometer,
		"switch":      NewSwitch,
		"flow_meter":  NewFlowMeter,
	}
	outputFactories = map[string]OutputDeviceFactory{
		"heater":     NewHeater,
		"cooler":     NewCooler,
		"thermostat": NewThermostat,
		"gpio":       NewGPIO,
		"recorder":   NewRecorder,
		"pwm":        NewPWM,
		"motor":      NewMotor,
		"file":       NewFile,
	}
)

func NewAppFactory() *AppFactory {
	a := &AppFactory{
		inputFactories:  inputFactories,
		outputFactories: outputFactories,
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
	f := FlowMeter{}
	th := Thermostat{}
	r := Recorder{}
	return map[string]ConfigHelper{
		"thermometer": t.Config(),
		"switch":      s.Config(),
		"gpio":        g.Config(),
		"heater":      h.Config(),
		"cooler":      c.Config(),
		"thermostat":  th.Config(),
		"recorder":    r.Config(),
		"flow_meter":  f.Config(),
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
	Status() interface{}
}

func deviceType(t string) string {
	_, ok := inputFactories[t]
	if ok {
		return "input"
	}
	_, ok = outputFactories[t]
	if ok {
		return "output"
	}
	return ""
}

func NewInputDevice(pin *Pin) (dev InputDevice, err error) {
	f, ok := inputFactories[pin.Type]
	if !ok {
		return nil, errors.New("invalid pin type")
	}
	return f(pin)
}

type OutputDeviceFactory func(pin *Pin) (OutputDevice, error)

//Outputdevices turn things on and off.  Currently the
type OutputDevice interface {
	On(val *Value) error
	Off() error
	Update(msg *Message)
	Status() interface{}
	Config() ConfigHelper
}

func NewOutputDevice(pin *Pin) (dev OutputDevice, err error) {
	f, ok := outputFactories[pin.Type]
	if !ok {
		return nil, errors.New("invalid pin type")
	}
	return f(pin)
}
