package output

import (
	"fmt"
	"bitbucket.com/cswank/gogadgets"
	"bitbucket.com/cswank/gogadgets/devices"
	"bitbucket.com/cswank/gogadgets/pins"
)

type Config struct {
	Location string
	Name string
	Pin *pins.Pin
}

type OutputGadget struct {
	gogadgets.Gadget
	Location string
	Name string
	Output devices.OutputDevice
	OnCommand string
	OffCommand string
	shutdown bool
	status bool
}


func NewOutputGadget(config *Config) (*OutputGadget, error) {
	dev, err := devices.NewOutputDevice(config.Pin)
	if err == nil {
		return &OutputGadget{
			Location: config.Location,
			Name: config.Name,
			OnCommand: fmt.Sprintf("turn on %s %s", config.Location, config.Name),
			OffCommand: fmt.Sprintf("turn off %s %s", config.Location, config.Name),
			Output: dev,
		}, err
	}
	return nil, err
}

func (od *OutputGadget) isMine(msg *gogadgets.Message) bool {
	return msg.Command == od.OnCommand || msg.Command == od.OffCommand || msg.Command == "shutdown"
}

func (od *OutputGadget) readCommand(msg *gogadgets.Message) {
	if msg.Command == "shutdown" {
		fmt.Println("shutting down")
		od.Output.Off()
		od.shutdown = true
	} else {
		
	}
}

func (od *OutputGadget) readMessage(msg *gogadgets.Message) {
	if msg.Type == "command" {
		od.readCommand(msg)
	}
}

func (od *OutputGadget) Start(input <-chan gogadgets.Message, output chan<- gogadgets.Message) {
	for !od.shutdown {
		msg := <-input
		if od.isMine(&msg) {
			od.readMessage(&msg)
		}
	}
}
