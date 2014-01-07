package gogadgets

import (
	"log"
	"fmt"
)

type App struct {
	gadgets []GoGadget
	masterHost string
	pubPort int
	subPort int
	channels map[string]chan Message
}

func NewApp(config *Config) *App {
	gadgets := GetGadgets(config.Gadgets)
	if  config.PubPort == 0 {
		config.SubPort = 6111
		config.PubPort = 6112
	}
	return &App{
		masterHost: config.MasterHost,
		pubPort: config.PubPort,
		subPort: config.SubPort,
		gadgets: gadgets,
	}
}

func GetGadgets(configs []GadgetConfig) []GoGadget {	
	g := make([]GoGadget, len(configs))
	for i, config := range configs {
		gadget, err := NewGadget(&config)
		if err != nil {
			log.Println(err)
		}
		g[i] = gadget
	}
	return g
}

func (a *App) Start(stop <-chan bool) {
	a.gadgets = append(a.gadgets, &Runner{})
	sockets := &Sockets{
		masterHost: a.masterHost,
		pubPort: a.pubPort,
		subPort: a.subPort,
	}
	a.gadgets = append(a.gadgets, sockets)
	in := make(chan Message, 100)
	a.channels = make(map[string]chan Message)
	for _, gadget := range a.gadgets {
		out := make(chan Message, 100)
		a.channels[gadget.GetUID()] = out
		go gadget.Start(out, in)
	}
	keepRunning := true
	log.Println("started gagdgets")
	for keepRunning {
		select {
		case msg := <-in:
			a.sendMessage(msg)
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

func (a *App) sendMessage(msg Message) {
	if msg.Target == "" {
		for uid, channel := range a.channels {
			if uid != msg.Sender {
				fmt.Println(uid)
				channel<- msg
				fmt.Println(uid)
			}
		}
	} else {
		channel, ok := a.channels[msg.Target]
		if ok {
			channel<- msg
		}
	}
}

func (a *App) AddGadget(gadget GoGadget) {
	a.gadgets = append(a.gadgets, gadget)
}
