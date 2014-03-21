package gogadgets

import (
	"time"
)

var (
	COMMAND      = "command"
	METHOD       = "method"
	DONE         = "done"
	UPDATE       = "update"
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
	if v.Units == "F" || v.Units == "f" {
		f = (f - 32.0) / 1.8
	} else if v.Units == "Gallons" || v.Units == "gallons" || v.Units == "gal" {
		f = f * 3.78541
	}
	return f, ok
}

type Info struct {
	Direction string `json:"direction"`
	On        string `json:"on"`
	Off       string `json:"off"`
}

type Method struct {
	Step  int      `json:"step"`
	Steps []string `json:"steps"`
	Time  int      `json:"time"`
}

//Message is what all Gadgets pass around to each
//other.
type Message struct {
	Name        string    `json:"name"`
	Location    string    `json:"location"`
	Type        string    `json:"type"`
	Sender      string    `json:"sender"`
	Target      string    `json:"target"`
	Body        string    `json:"body"`
	Method      Method    `json:"method"`
	Timestamp   time.Time `json:"timestamp"`
	Value       Value     `json:"value"`
	TargetValue *Value    `json:"targetValue"`
	Info        Info      `json:"info"`
}

type Pin struct {
	Type      string  `json:"type"`
	Port      string  `json:"port"`
	Pin       string  `json:"pin"`
	Direction string  `json:"direction"`
	Edge      string  `json:"edge"`
	OneWireId string  `json:"onewireId"`
	Value     interface{}  `json:"value"`
	Units     string  `json:"units"`
	Args      map[string]string  `json:"args"`
}

type GadgetConfig struct {
	Location   string `json:"location"`
	Name       string `json:"name"`
	OnCommand  string `json:"onCommand"`
	OffCommand string `json:"offCommand"`
	Pin        Pin    `json:"pin"`
}

type Config struct {
	Host    string         `json:"host"`
	PubPort int            `json:"pubPort"`
	SubPort int            `json:"subPort"`
	Gadgets []GadgetConfig `json:"gadgets"`
}
