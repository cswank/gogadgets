package models

import (
	"time"
)

var (
	COMMAND      = "command"
	METHOD       = "method"
	DONE         = "done"
	UPDATE       = "update"
	GADGET       = "gadget"
	METHODUPDATE = "method update"
)

type GoGadget interface {
	GetUID() string
	Start(input <-chan Message, output chan<- Message)
}

type Value struct {
	Value interface{} `json:"value,omitempty"`
	Units string      `json:"units,omitempty"`
	ID    string      `json:"id,omitempty"`
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
	Direction string `json:"direction,omitempty"`
	On        string `json:"on,omitempty"`
	Off       string `json:"off,omitempty"`
}

type Method struct {
	Step  int      `json:"step,omitempty"`
	Steps []string `json:"steps,omitempty"`
	Time  int      `json:"time,omitempty"`
}

//Message is what all Gadgets pass around to each
//other.
type Message struct {
	Name        string       `json:"name,omitempty"`
	Location    string       `json:"location,omitempty"`
	Type        string       `json:"type,omitempty"`
	Sender      string       `json:"sender,omitempty"`
	Target      string       `json:"target,omitempty"`
	Body        string       `json:"body,omitempty"`
	Method      Method       `json:"method,omitempty"`
	Timestamp   time.Time    `json:"timestamp,omitempty"`
	Value       Value        `json:"value,omitempty"`
	TargetValue *Value       `json:"targetValue,omitempty"`
	Info        Info         `json:"info,omitempty"`
	Config      Config       `json:"config,omitempty"`
}

type Pin struct {
	Type      string                 `json:"type,omitempty"`
	Port      string                 `json:"port,omitempty"`
	Pin       string                 `json:"pin,omitempty"`
	Direction string                 `json:"direction,omitempty"`
	Edge      string                 `json:"edge,omitempty"`
	OneWireId string                 `json:"onewireId,omitempty"`
	Value     interface{}            `json:"value,omitempty"`
	Units     string                 `json:"units,omitempty"`
	Platform  string                 `json:"platform,omitempty"`
	Frequency int                    `json:"frequency,omitempty"`
	Args      map[string]interface{} `json:"args,omitempty"`
	Pins      map[string]Pin         `json:"pins,omitempty"`
}

type GadgetConfig struct {
	Location     string `json:"location,omitempty"`
	Name         string `json:"name,omitempty"`
	OnCommand    string `json:"onCommand,omitempty"`
	OffCommand   string `json:"offCommand,omitempty"`
	InitialValue string `json:"initialValue,omitempty"`
	Pin          Pin    `json:"pin,omitempty"`
}

type Config struct {
	Host    string         `json:"host,omitempty"`
	PubPort int            `json:"pubPort,omitempty"`
	SubPort int            `json:"subPort,omitempty"`
	Gadgets []GadgetConfig `json:"gadgets,omitempty"`
}

type ConfigHelper struct {
	Fields map[string][]string `json:"fields"`
	Units []string `json:"units,omitempty"`
	Args map[string]interface{} `json:"args,omitempty"`
	Pins map[string]map[string]string `json:"pins,omitempty"`
	PinType string `json:"pinType"`
}
