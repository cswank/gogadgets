package gogadgets

import (
	"errors"
	"fmt"
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

//There are 5 types of Input/Output devices build into
//GoGadgets (header, cooler, gpio, thermometer and switch)
//NewGadget reads a GadgetConfig and creates the correct
//type of Gadget.
func NewGadget(config *GadgetConfig) (Gadgeter, error) {
	if config.Type == "cron" {
		return newSystemGadget(config)
	}
	switch deviceType(config.Pin.Type) {
	case "input":
		return NewInputGadget(config)
	case "output":
		return NewOutputGadget(config)
	}
	return nil, fmt.Errorf(
		"couldn't build a gadget based on config: %s %s",
		config.Location,
		config.Name,
	)
}

func newSystemGadget(config *GadgetConfig) (Gadgeter, error) {
	if config.Type == "cron" {
		return NewCron(config)
	}
	return nil, fmt.Errorf("don't know how to build %s", config.Name)
}

//Input Gadgets read from input devices and report their values (thermometer
//is an example).
func NewInputGadget(config *GadgetConfig) (gadget *Gadget, err error) {
	dev, err := NewInputDevice(&config.Pin)
	if err == nil {
		gadget = &Gadget{
			Location:   config.Location,
			Name:       config.Name,
			Input:      dev,
			Direction:  "input",
			OnCommand:  "n/a",
			OffCommand: "n/a",
			UID:        fmt.Sprintf("%s %s", config.Location, config.Name),
		}
	}
	return gadget, err
}

//Output Gadgets turn devices on and off.
func NewOutputGadget(config *GadgetConfig) (gadget *Gadget, err error) {
	dev, err := NewOutputDevice(&config.Pin)
	if config.OnCommand == "" {
		config.OnCommand = fmt.Sprintf("turn on %s %s", config.Location, config.Name)
	}
	if config.OffCommand == "" {
		config.OffCommand = fmt.Sprintf("turn off %s %s", config.Location, config.Name)
	}
	if err == nil {
		gadget = &Gadget{
			Location:       config.Location,
			Name:           config.Name,
			Direction:      "output",
			OnCommand:      config.OnCommand,
			OffCommand:     config.OffCommand,
			InitialValue:   config.InitialValue,
			Output:         dev,
			Operator:       ">=",
			UID:            fmt.Sprintf("%s %s", config.Location, config.Name),
			filterMessages: config.Pin.Type != "recorder",
		}
	} else {
		panic(err)
	}
	return gadget, err
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
