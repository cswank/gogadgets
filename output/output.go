package output

import (
	"fmt"
	"strings"
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
	uid string
	shutdown bool
	status bool
	units string
	trigger *Trigger
	out chan<- gogadgets.Message
	triggerIn <-chan gogadgets.Message
	triggerOut chan<- gogadgets.Message
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
			uid: fmt.Sprintf("%s %s", config.Location, config.Name),
		}, err
	}
	return nil, err
}

func (od *OutputGadget) isMine(msg *gogadgets.Message) bool {
	return msg.Body == od.OnCommand || msg.Body == od.OffCommand || msg.Body == "shutdown"
}

func (od *OutputGadget) Start(in <-chan gogadgets.Message, out chan<- gogadgets.Message) {
	od.out = out
	od.triggerIn = make(chan gogadgets.Message)
	od.triggerOut = make(chan gogadgets.Message)
	for !od.shutdown {
		select {
		case msg := <-in:
			if od.isMine(&msg) {
				od.readMessage(&msg)
			}
			if od.trigger != nil {
				od.triggerOut<- msg
			}
		case msg := <-od.triggerIn:
			od.readTrigger(&msg)
		}
	}
}

func (od *OutputGadget) readCommand(msg *gogadgets.Message) {
	if msg.Body == "shutdown" {
		od.Output.Off()
		od.shutdown = true
	} else if strings.Index(msg.Body, od.OnCommand) == 0 {
		od.readOnCommand(msg)
	} else if strings.Index(msg.Body, od.OffCommand) == 0 {
		od.readOffCommand(msg)
	}
}

func (od *OutputGadget) readTrigger(msg *gogadgets.Message) {
	if msg.Type == gogadgets.DONE || msg.Type == gogadgets.UPDATE {
		od.out<- *msg
	} else if msg.Type == gogadgets.COMMAND {
		od.readCommand(msg)
	}
}

func (od *OutputGadget) readMessage(msg *gogadgets.Message) {
	if msg.Type == gogadgets.COMMAND {
		od.readCommand(msg)
	}
	od.sendStatus()
}

func (od *OutputGadget) readOnCommand(msg *gogadgets.Message) {
	od.status = true
	od.Output.On()
}

func (od *OutputGadget) readOffCommand(msg *gogadgets.Message) {
	od.status = false
	od.Output.Off()
}

func (od *OutputGadget) sendStatus() {
	location := gogadgets.Location{
		Output: map[string]gogadgets.Device{
			od.Name: gogadgets.Device{
				Units: od.units,
				Value:od.Output.Status(),
			},
		},
	}
	msg := gogadgets.Message{
		Sender: od.uid,
		Type: gogadgets.STATUS,
		Locations: map[string]gogadgets.Location{od.Location: location},
	}
	od.out<- msg
}
