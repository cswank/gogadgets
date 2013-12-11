package output

import (
	"fmt"
	"time"
	"bitbucket.com/cswank/gogadgets/io"
	"bitbucket.com/cswank/gogadgets/pins"
)

type OutputGadget struct {
	Gadget
	Location string
	Name string
	Output OutputDevice
	On string
	Off string
	exit bool
	status bool
}

type GPIO struct {
	OutputDevice
	Pin: *pins.GPIO,
}

func NewOutputDevice(pin *pins.GPIO) OutputDevice {
	
}

func NewOutputGadget(config *Config) *OutputGadget {
	return &OutputGadget{
		Location: config.Location,
		Name: config.Name,
		On: fmt.Sprintf("turn on %s %s", config.Location, config.Name),
		Off: fmt.Sprintf("turn off %s %s", config.Location, config.Name),
		Pin: NewOutputDevice(config.Pin),
	}
}

func (od *OutputGadget) isMine(msg *Message) bool {
	return msg.Command == od.On || msg.Command == od.Off
}

func (od *OutputGadget) readCommand(msg *Message) {
	if msg.Command == "exit" {
		od.exit = true
	} else {
		
	}
}

func (od *OutputGadget) readMessage(msg *Message) {
	if msg.Type == "command" {
		od.readCommand(msg)
	}
}

func (od *OutputGadget) Start(input <-chan Message, output chan<- Message) {
	for !od.exit {
		msg := <-input
		if od.isMine(&msg) {
			od.readMessage(&msg)
		}
	}
}
