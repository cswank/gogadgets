package gogadgets

import (
	"time"
)

var (
	COMMAND = "command"
	DONE = "done"
	UPDATE = "update"
	STATUS = "status"
)

type Gadget interface {
	Start(input <-chan Message, output chan<- Message)
}

type Device struct {
	Units string      `json:"units"`
	Value interface{} `json:"value"`
	ID    string      `json:"id"`
}

type Location struct {
	Input  map[string]Device `json:"input"`
	Output map[string]Device `json:"output"`
}

type Message struct {
	Sender      string      `json:"sender"`
	Type        string      `json:"type"`
	Body        string      `json:"body"`
	Timestamp   time.Time   `json:"timestamp"`
	Name        string      `json:"name"`
	Locations   map[string]Location `json:"locations"`
}
