package main

import (
	"fmt"
	"time"
	"bitbucket.com/cswank/gogadgets"
)

type Greenhouse struct {
	gogadgets.GoGadget
	temperature float64
	sleepTimes map[string]time.Duration
	out chan<- gogadgets.Message
}

func (g *Greenhouse) getMessage(cmd, location string) gogadgets.Message {
	return gogadgets.Message{
		Sender: "greenhouse watcher",
		Type: "command",
		Body: fmt.Sprintf("turn %s %s pump", cmd, location),
	}
}

func (g *Greenhouse) wait(location string) {
	time.Sleep(g.sleepTimes[location])
	offCmd := g.getMessage("on", location)
	g.out<- offCmd
}

func (g *Greenhouse) Start(in <-chan gogadgets.Message, out chan<- gogadgets.Message) {
	g.out = out
	for {
		msg := <-in
		if msg.Type == "update" &&
			msg.Location == "greenhouse" &&
			msg.Name == "temperature" {
			g.temperature = msg.Value.Value.(float64)
		} else if msg.Type == "update" &&
			msg.Name == "switch" &&
			msg.Value.Value == false {
			cmd := g.getMessage("off", msg.Location)
			out<- cmd
			if g.temperature >= 12.0 {
				go g.wait(msg.Location)
			}
		} else if msg.Type == "command" && msg.Body == "shutdown" {
			out<- gogadgets.Message{}
			return
		}
	}
}

func main() {
	a := gogadgets.App{}
	g := &Greenhouse{}
	a.AddGadget(g)
	fmt.Println(a)
}
