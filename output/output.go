package output

import (
	"fmt"
	"strings"
	"bitbucket.com/cswank/gogadgets"
	"bitbucket.com/cswank/gogadgets/devices"
)

type Config struct {
	Location string
	Name string
	Pin *devices.Pin
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
	triggerIn chan gogadgets.Message
	triggerOut chan gogadgets.Message
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

func (og *OutputGadget) isMine(msg *gogadgets.Message) bool {
	return strings.Index(msg.Body, og.OnCommand) == 0 || strings.Index(msg.Body, og.OffCommand) == 0 || msg.Body == "shutdown"
}

func (og *OutputGadget) Start(in <-chan gogadgets.Message, out chan<- gogadgets.Message) {
	og.out = out
	og.triggerIn = make(chan gogadgets.Message)
	og.triggerOut = make(chan gogadgets.Message)
	for !og.shutdown {
		select {
		case msg := <-in:
			if og.isMine(&msg) {
				og.readMessage(&msg)
			}
			if og.trigger != nil && msg.Type == gogadgets.STATUS { 
				og.triggerOut<- msg
			}
		case msg := <-og.triggerIn:
			og.readTrigger(&msg)
		}
	}
}

func (og *OutputGadget) readTrigger(msg *gogadgets.Message) {
	if msg.Type == gogadgets.DONE || msg.Type == gogadgets.UPDATE {
		og.out<- *msg
	} else if msg.Type == gogadgets.COMMAND {
		og.readCommand(msg)
	}
}

func (og *OutputGadget) readMessage(msg *gogadgets.Message) {
	if msg.Type == gogadgets.COMMAND {
		og.readCommand(msg)
	}
	
}

func (og *OutputGadget) readCommand(msg *gogadgets.Message) {
	if msg.Body == "shutdown" {
		og.Output.Off()
		og.shutdown = true
		og.sendStatus()
	} else if strings.Index(msg.Body, og.OnCommand) == 0 {
		og.readOnCommand(msg)
	} else if strings.Index(msg.Body, og.OffCommand) == 0 {
		og.readOffCommand(msg)
	}
}

func (og *OutputGadget) readOnCommand(msg *gogadgets.Message) {
	if !og.status {
		og.status = true
		og.Output.On()
		if len(msg.Body) > len(og.OnCommand) {
			og.startTrigger(msg)
		}
		og.sendStatus()
	}
}

func (og *OutputGadget) startTrigger(msg *gogadgets.Message) {
	og.trigger = &Trigger{
		location: og.Location,
		name: og.Name,
		operator: ">=",
		command: msg.Body[len(og.OnCommand):],
		offCommand: og.OffCommand,
	}
	go og.trigger.Start(og.triggerIn, og.triggerOut)
}

func (og *OutputGadget) readOffCommand(msg *gogadgets.Message) {
	if og.status {
		if og.trigger != nil && msg.Sender != og.trigger.uid {
			og.triggerOut<- gogadgets.Message{
				Type: gogadgets.COMMAND,
				Body: "stop",
			}
			og.trigger = nil
		}
		og.status = false
		og.Output.Off()
		og.sendStatus()
	}
}

func (og *OutputGadget) sendStatus() {
	location := gogadgets.Location{
		Output: map[string]gogadgets.Device{
			og.Name: gogadgets.Device{
				Units: og.units,
				Value:og.Output.Status(),
			},
		},
	}
	msg := gogadgets.Message{
		Sender: og.uid,
		Type: gogadgets.STATUS,
		Locations: map[string]gogadgets.Location{og.Location: location},
	}
	og.out<- msg
}
