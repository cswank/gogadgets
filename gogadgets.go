package gogadgets

import (
	"errors"
	"fmt"
)

var (
	inputs = map[string]CreateInputDevice{
		"thermometer": NewThermometer,
		"switch":      NewSwitch,
		"flow_meter":  NewFlowMeter,
		"xbee":        NewXBee,
	}

	outputs = map[string]CreateOutputDevice{
		"alarm":      NewAlarm,
		"heater":     NewHeater,
		"cooler":     NewCooler,
		"thermostat": NewThermostat,
		"boiler":     NewBoiler,
		"gpio":       NewGPIO,
		"recorder":   NewRecorder,
		"pwm":        NewPWM,
		"motor":      NewMotor,
		"file":       NewFile,
		"sms":        NewSMS,
	}
)

// NewGadget returns a gadget with an input or output device
// There are several types of Input/Output devices build into
// GoGadgets (eg: header, cooler, gpio, thermometer and switch)
// NewGadget reads a GadgetConfig and creates the correct
// type of Gadget.
func NewGadget(config *GadgetConfig) (Gadgeter, error) {
	if config.Type == "cron" {
		return newSystemGadget(config)
	}
	switch deviceType(config.Pin.Type) {
	case "input":
		return newInputGadget(config)
	case "output":
		return newOutputGadget(config)
	default:
		return nil, fmt.Errorf("couldn't build a gadget based on config: %v", config)
	}
}

func newSystemGadget(config *GadgetConfig) (Gadgeter, error) {
	if config.Type == "cron" {
		return NewCron(config)
	}
	return nil, fmt.Errorf("don't know how to build %s", config.Name)
}

// InputGadgets read from input devices and report their values (thermometer
// is an example).
func newInputGadget(config *GadgetConfig) (gadget *Gadget, err error) {
	dev, err := NewInputDevice(&config.Pin)
	m := map[bool]string{}
	if config.OnValue != "" {
		m[true] = config.OnValue
	}

	if config.OffValue != "" {
		m[false] = config.OffValue
	}

	if err == nil {
		gadget = &Gadget{
			Type:        config.Pin.Type,
			Location:    config.Location,
			Name:        config.Name,
			Input:       dev,
			Direction:   "input",
			UID:         fmt.Sprintf("%s %s", config.Location, config.Name),
			inputLookup: m,
		}
	}
	return gadget, err
}

//Output Gadgets turn devices on and off.
func newOutputGadget(config *GadgetConfig) (gadget *Gadget, err error) {
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
		filterMessages: config.Pin.Type != "recorder" && config.Pin.Type != "alarm",
	}
	return gadget, nil
}

func RegisterInput(name string, f CreateInputDevice) {
	inputs[name] = f
}

func RegisterOutput(name string, f CreateOutputDevice) {
	outputs[name] = f
}

type CreateInputDevice func(pin *Pin, opts ...func(InputDevice) error) (InputDevice, error)

//Inputdevices are started as goroutines by the Gadget
//that contains it.
type InputDevice interface {
	Start(<-chan Message, chan<- Value)
	GetValue() *Value
	Config() ConfigHelper
}

type Poller interface {
	Wait() error
	Status() map[string]bool
}

func deviceType(t string) string {
	_, ok := inputs[t]
	if ok {
		return "input"
	}
	_, ok = outputs[t]
	if ok {
		return "output"
	}
	return ""
}

func NewInputDevice(pin *Pin) (dev InputDevice, err error) {
	f, ok := inputs[pin.Type]
	if !ok {
		return nil, errors.New("invalid pin type")
	}
	return f(pin)
}

type CreateOutputDevice func(pin *Pin) (OutputDevice, error)

type Commands struct {
	On  []string
	Off []string
}

// OutputDevice turns things on and off.  Currently the
type OutputDevice interface {
	On(val *Value) error
	Off() error
	Update(msg *Message) bool
	Status() map[string]bool
	Config() ConfigHelper
	Commands(string, string) *Commands
}

func NewOutputDevice(pin *Pin) (dev OutputDevice, err error) {
	f, ok := outputs[pin.Type]
	if !ok {
		return nil, errors.New("invalid pin type")
	}

	for k, p := range pin.Pins {
		f, ok := outputs[p.Type]
		if !ok {
			return nil, errors.New("invalid pin type")
		}
		p.new = f
		pin.Pins[k] = p
	}

	return f(pin)
}
