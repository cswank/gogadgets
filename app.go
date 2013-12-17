package main

import (
	"bitbucket.com/cswank/gogadgets/models"
	"bitbucket.com/cswank/gogadgets/gadgets"
)

type Gadgets struct {
	configs []models.Config
	gadgets []gadgets.Gadget
}

type GadgetHolder struct {
	gadget models.Gadget
	out chan models.Message
}

func (g *Gadgets) Start(stop <-chan bool) {
	in := make(chan models.Message)
	gadgets := []GadgetHolder{}
	for _, config := range g.configs {
		out := make(chan models.Message)
		gadget := NewGadet(config)
		go gadget.Start(out, in)
		gadgets = append(gadgets, GadgetHolder{gadget, out})
	}
	keepRunning := true
	for keepRunning {
		select {
		case msg := <-in:
			for _, g := range gadgets {
				g.out<- msg
			}
		case keepRunning = <-stop:
			stopMessage := models.Message{
				Type: models.COMMAND,
				Body: "shutdown",
			}
			for _, g := range gadgets {
				g.out<- stopMessage
			}
			for _, g := range gadgets {
				<-in
			}
		}
	}
}

func main() {
	
}

