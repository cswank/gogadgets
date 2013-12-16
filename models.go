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

type Value struct {
	Value interface{} `json:"value"`
	Units string      `json:"units"`
	ID    string      `json:"id"`
}

type Message struct {
	Sender      string      `json:"sender"`
	Type        string      `json:"type"`
	Body        string      `json:"body"`
	Timestamp   time.Time   `json:"timestamp"`
	Name        string      `json:"name"`
	Location    string      `json:"location"`
	Value       Value       `json:"value"`
}
