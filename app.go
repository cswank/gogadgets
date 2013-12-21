package gogadgets

import (
	"log"
)

type App struct {
	gadgets []GoGadget
	channels []chan Message
}

func (a *App) Start(stop <-chan bool) {
	in := make(chan Message)
	n := len(a.gadgets) + 1
	a.channels = make([]chan Message, n)
	a.channels[0] = make(chan Message)
	runner := Runner{}
	go runner.Start(a.channels[0], in)
	for i, gadget := range a.gadgets {
		out := make(chan Message)
		go gadget.Start(out, in)
		a.channels[i + 1] = out
	}
	keepRunning := true
	for keepRunning {
		select {
		case msg := <-in:
			for _, channel := range a.channels {
				channel<- msg
			}
		case keepRunning = <-stop:
			stopMessage := Message{
				Type: COMMAND,
				Body: "shutdown",
			}
			for _, channel := range a.channels {
				channel<- stopMessage
			}
			for _, _  = range a.channels {
				<-in
			}
		}
	}
}

func (a *App) AddGadget(gadget GoGadget) {
	a.gadgets = append(a.gadgets, gadget)
}

func GatGadgets(configs []*Config) *App {
	g := make([]GoGadget, len(configs))
	for i, config := range configs {
		gadget, err := NewGadget(config)
		if err != nil {
			log.Println(err)
		}
		g[i] = gadget
	}
	return &App{gadgets:g}
}
