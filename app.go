package gogadgets

import (
	"log"
)

type App struct {
	gadgets []GoGadget
	MasterHost string
	channels map[string]chan Message
}

func (a *App) Start(stop <-chan bool) {
	a.gadgets = append(a.gadgets, &Runner{})
	//a.gadgets = append(a.gadgets, &Sockets{masterHost: a.MasterHost})
	in := make(chan Message)
	a.channels = make(map[string]chan Message)
	for _, gadget := range a.gadgets {
		out := make(chan Message)
		a.channels[gadget.GetUID()] = out
		go gadget.Start(out, in)
	}
	keepRunning := true
	for keepRunning {
		select {
		case msg := <-in:
			for uid, channel := range a.channels {
				if uid != msg.Sender {
					channel<- msg
				}
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
