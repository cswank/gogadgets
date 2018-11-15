package gogadgets

import (
	"errors"
	"fmt"
	"sync"
	"time"
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

var (
	COMMAND      = "command"
	ERROR        = "error"
	METHOD       = "method"
	DONE         = "done"
	UPDATE       = "update"
	GADGET       = "gadget"
	STATUS       = "status"
	METHODUPDATE = "method update"
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

type GoGadgeter interface {
	GetUID() string
	Start(input <-chan Message, output chan<- Message)
}

type Value struct {
	Value    interface{}     `json:"value,omitempty"`
	Units    string          `json:"units,omitempty"`
	Output   map[string]bool `json:"io,omitempty"`
	ID       string          `json:"id,omitempty"`
	Cmd      string          `json:"command,omitempty"`
	location string
	name     string
}

func (v *Value) GetName() string {
	return v.name
}

func (v *Value) ToFloat() (f float64, ok bool) {
	switch V := v.Value.(type) {
	case bool:
		if V {
			f = 1.0
		} else {
			f = 0.0
		}
		ok = true
	case float64:
		f = V
		ok = true
	}
	return f, ok
}

type Info struct {
	Direction string   `json:"direction,omitempty"`
	Type      string   `json:"type,omitempty"`
	On        []string `json:"on,omitempty"`
	Off       []string `json:"off,omitempty"`
}

type Method struct {
	Step  int      `json:"step,omitempty"`
	Steps []string `json:"steps,omitempty"`
	Time  int      `json:"time,omitempty"`
}

//Message is what all Gadgets pass around to each
//other.
type Message struct {
	UUID        string    `json:"uuid"`
	From        string    `json:"from,omitempty"`
	Name        string    `json:"name,omitempty"`
	Location    string    `json:"location,omitempty"`
	Type        string    `json:"type,omitempty"`
	Sender      string    `json:"sender,omitempty"`
	Target      string    `json:"target,omitempty"`
	Body        string    `json:"body,omitempty"`
	Host        string    `json:"host,omitempty"`
	Method      Method    `json:"method,omitempty"`
	Timestamp   time.Time `json:"timestamp,omitempty"`
	Value       Value     `json:"value,omitempty"`
	TargetValue *Value    `json:"target_value,omitempty"`
	Info        Info      `json:"info,omitempty"`
	Config      Config    `json:"config,omitempty"`
}

type Pin struct {
	Type        string        `json:"type,omitempty"`
	Port        string        `json:"port,omitempty"`
	Pin         string        `json:"pin,omitempty"`
	Direction   string        `json:"direction,omitempty"`
	Edge        string        `json:"edge,omitempty"`
	ActiveLow   string        `json:"active_low,omitempty"`
	OneWirePath string        `json:"onewire_path,omitempty"`
	OneWireId   string        `json:"onewire_id,omitempty"`
	Sleep       time.Duration `json:"sleep,omitempty"`
	Value       interface{}   `json:"value,omitempty"`
	Units       string        `json:"units,omitempty"`
	//Platform is either "rpi" or "beaglebone"
	Platform  string                 `json:"platform,omitempty"`
	Frequency int                    `json:"frequency,omitempty"`
	Args      map[string]interface{} `json:"args,omitempty"`
	Pins      map[string]Pin         `json:"pins,omitempty"`
	Lock      *sync.Mutex            `json:"-"`

	new func(*Pin) (OutputDevice, error)
}
