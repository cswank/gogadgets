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
		"boiler":     NewBoiler,
		"gpio":       GPIOFactory,
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
			Type:      config.Pin.Type,
			Location:  config.Location,
			Name:      config.Name,
			Input:     dev,
			Direction: "input",
			UID:       fmt.Sprintf("%s %s", config.Location, config.Name),
		}
	}
	return gadget, err
}

//Output Gadgets turn devices on and off.
func NewOutputGadget(config *GadgetConfig) (gadget *Gadget, err error) {
	dev, err := NewOutputDevice(&config.Pin)
	if err != nil {
		panic(err)
	}

	if cmds := dev.Commands(config.Location, config.Name); cmds != nil {
		config.OnCommands = cmds.On
		config.OffCommands = cmds.Off
	} else {
		if len(config.OnCommands) == 0 {
			config.OnCommands = []string{fmt.Sprintf("turn on %s %s", config.Location, config.Name)}
		}
		if len(config.OffCommands) == 0 {
			config.OffCommands = []string{fmt.Sprintf("turn off %s %s", config.Location, config.Name)}
		}
	}
	gadget = &Gadget{
		Type:           config.Pin.Type,
		Location:       config.Location,
		Name:           config.Name,
		Direction:      "output",
		OnCommands:     config.OnCommands,
		OffCommands:    config.OffCommands,
		InitialValue:   config.InitialValue,
		Output:         dev,
		Operator:       ">=",
		UID:            fmt.Sprintf("%s %s", config.Location, config.Name),
		filterMessages: config.Pin.Type != "recorder",
	}
	return gadget, nil
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
	Status() map[string]bool
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

type Commands struct {
	On  []string
	Off []string
}

//Outputdevices turn things on and off.  Currently the
type OutputDevice interface {
	On(val *Value) error
	Off() error
	Update(msg *Message) bool
	Status() map[string]bool
	Config() ConfigHelper
	Commands(string, string) *Commands
}

func NewOutputDevice(pin *Pin) (dev OutputDevice, err error) {
	f, ok := outputFactories[pin.Type]
	if !ok {
		return nil, errors.New("invalid pin type")
	}
	return f(pin)
}
