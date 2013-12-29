package main

import (
	"fmt"
	"time"
	"io/ioutil"
	"encoding/json"
	"bitbucket.com/cswank/gogadgets"
)

type Greenhouse struct {
	gogadgets.GoGadget
	temperature float64
	sleepTimes map[string]time.Duration
	out chan<- gogadgets.Message
	status bool
}

func (g *Greenhouse)getMessage(cmd, location string) gogadgets.Message {
	return gogadgets.Message{
		Sender: "greenhouse watcher",
		Type: "command",
		Body: fmt.Sprintf("turn %s %s pump", cmd, location),
	}
}

func (g *Greenhouse)wait(location string) {
	time.Sleep(g.sleepTimes[location])
	offCmd := g.getMessage("on", location)
	g.out<- offCmd
}

func (g *Greenhouse)startPumps(location string) {
	for key, _ := range g.sleepTimes {
		msg := g.getMessage("on", key)
		g.out<- msg
	}
}

func (g *Greenhouse)Start(in <-chan gogadgets.Message, out chan<- gogadgets.Message) {
	g.out = out
	for {
		msg := <-in
		if msg.Type == "update" &&
			msg.Location == "greenhouse" &&
			msg.Name == "temperature" {
			g.temperature = msg.Value.Value.(float64)
			if g.temperature >= 12.0 && !g.status {
				g.status = true
				g.startPumps()
			}
		} else if msg.Type == "update" &&
			msg.Name == "switch" &&
			msg.Value.Value == false {
			cmd := g.getMessage("off", msg.Location)
			out<- cmd
			if g.temperature >= 12.0 {
				go g.wait(msg.Location)
			} else {
				g.status = false
			}
		} else if msg.Type == "command" && msg.Body == "shutdown" {
			out<- gogadgets.Message{}
			return
		}
	}
}

func main() {
	b, err := ioutil.ReadFile("config.json")
	if err != nil {
		panic(err)
	}
	cfg := &gogadgets.Config{}
	err = json.Unmarshal(b, cfg)
	a := gogadgets.NewApp(cfg)
	g := &Greenhouse{}
	a.AddGadget(g)
	stop := make(chan bool)
	a.Start(stop)
}
