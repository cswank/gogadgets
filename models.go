package gogadgets

import (
	"time"
)

var (
	COMMAND = "command"
	METHOD = "method"
	DONE = "done"
	UPDATE = "update"
	STATUS = "status"
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

type Info struct {
	Direction string      `json:"direction"`
	On string             `json:"on"`
	Off string            `json:"off"`
}

type Message struct {
	Sender      string      `json:"sender"`
	Type        string      `json:"type"`
	Body        string      `json:"body"`
	Method      []string    `json:"method"`
	Timestamp   time.Time   `json:"timestamp"`
	Name        string      `json:"name"`
	Location    string      `json:"location"`
	Value       Value       `json:"value"`
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
	Pin Pin
}

type Config struct {
	MasterHost string
	PubPort int
	SubPort int
	Gadgets []GadgetConfig
}

