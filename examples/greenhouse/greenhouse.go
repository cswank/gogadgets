package main

import (
	//"syscall"
	"flag"
	"fmt"
	"github.com/cswank/gogadgets"
	"github.com/cswank/gogadgets/utils"
	"io/ioutil"
	"time"
)

var (
	configFlag = flag.String("c", "", "Path to the config json file")
)

type Greenhouse struct {
	gogadgets.GoGadget
	temperature float64
	sleepTimes  map[string]time.Duration
	out         chan<- gogadgets.Message
	status      bool
	pumps       map[string]bool
}

func (g *Greenhouse) getMessage(cmd, location string) gogadgets.Message {
	return gogadgets.Message{
		Sender: "greenhouse watcher",
		Type:   "command",
		Body:   fmt.Sprintf("turn %s %s pump", cmd, location),
	}
}

func (g *Greenhouse) wait(location string) {
	time.Sleep(g.sleepTimes[location])
	cmd := g.getMessage("on", location)
	g.pumps[location] = true
	g.out <- cmd
}

func (g *Greenhouse) GetUID() string {
	return "greenhouse watcher"
}

func (g *Greenhouse) startPump(name string) {
	msg := g.getMessage("on", name)
	g.out <- msg
}

func (g *Greenhouse) startPumps() {
	for name, _ := range g.sleepTimes {
		g.pumps[name] = true
		go g.startPump(name)
	}
}

func (g *Greenhouse) Start(in <-chan gogadgets.Message, out chan<- gogadgets.Message) {
	g.out = out
	g.pumps = make(map[string]bool)
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
			g.pumps[msg.Location] &&
			msg.Value.Value.(float64) == 0.0 {
			cmd := g.getMessage("off", msg.Location)
			g.pumps[msg.Location] = false
			go func() {
				out <- cmd
			}()
			if g.temperature >= 12.0 {
				go g.wait(msg.Location)
			} else {
				g.status = false
			}
		} else if msg.Type == "command" && msg.Body == "shutdown" {
			go func() {
				out <- gogadgets.Message{}
			}()
			return
		}
	}
}

func main() {
	flag.Parse()
	if !utils.FileExists("/sys/bus/w1/devices/28-0000047ade8f") {
		ioutil.WriteFile("/sys/devices/bone_capemgr.9/slots", []byte("BB-W1:00A0"), 0666)
	}
	a := gogadgets.NewApp(configFlag)
	sleepTimes := map[string]time.Duration{
		"bed 1": 300 * time.Second,
		"bed 2": 300 * time.Second,
		"bed 3": 300 * time.Second,
	}
	g := &Greenhouse{sleepTimes: sleepTimes}
	a.AddGadget(g)
	a.Start()
}
