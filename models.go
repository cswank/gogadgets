package gogadgets

import (
	"time"
)

var (
	COMMAND = "command"
	METHOD = "method"
	DONE = "done"
	UPDATE = "update"
	METHODUPDATE = "method update"
)

type InputDeviceFactory func(pin *Pin) (InputDevice, error)

type OutputDeviceFactory func(pin *Pin) (OutputDevice, error)

type GoGadget interface {
	GetUID() string
	Start(input <-chan Message, output chan<- Message)
}

type Value struct {
	Value interface{} `json:"value"`
	Units string      `json:"units"`
	ID    string      `json:"id"`
}

func (v *Value) ToFloat() (f float64, ok bool) {
	f, ok = v.Value.(float64)
	if v.Units == "F" || v.Units == "F" {
		f = (f - 32.0) / 1.8
	} else if v.Units == "Gallons" || v.Units == "gallons" || v.Units == "gal" {
		f = f * 3.78541
	}
	return f, ok
}

type Info struct {
	Direction string      `json:"direction"`
	On string             `json:"on"`
	Off string            `json:"off"`
}

type Method struct {
	Step int        `json:"step"`
	Steps []string  `json:"steps"`
	Time int        `json:"time"`
}

type Message struct {
	Sender      string      `json:"sender"`
	Target      string      `json:"target"`
	Type        string      `json:"type"`
	Body        string      `json:"body"`
	Method      Method      `json:"method"`
	Timestamp   time.Time   `json:"timestamp"`
	Name        string      `json:"name"`
	Location    string      `json:"location"`
	Value       Value       `json:"value"`
	TargetValue *Value       `json:"targetValue"`
	Info        Info        `json:"info"`
}

type Pin struct {
	Type string
	Port string
	Pin string
	Direction string
	Edge string
	OneWireId string
	Value interface{}
	Units string
}

type GadgetConfig struct {
	Location string
	Name string
	OnCommand string
	OffCommand string
	Pin Pin
}

type Config struct {
	MasterHost string
	PubPort int
	SubPort int
	Gadgets []GadgetConfig
}

