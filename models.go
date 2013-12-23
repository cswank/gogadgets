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

type Message struct {
	Sender      string      `json:"sender"`
	Type        string      `json:"type"`
	Body        string      `json:"body"`
	Method      []string    `json:"method"`
	Timestamp   time.Time   `json:"timestamp"`
	Name        string      `json:"name"`
	Location    string      `json:"location"`
	Value       Value       `json:"value"`
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

type Config struct {
	Location string
	Name string
	Pin Pin
}
