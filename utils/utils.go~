package gogadgets

import (
	"fmt"
	"time"
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
	Command     string      `json:"command"`
	Timestamp   time.Time   `json:"timestamp"`
	Name        string      `json:"name"`
	Locations   map[string]Location `json:"locations"`
}
