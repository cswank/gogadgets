package gogadgets

import (
	"log"
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
		log.Println(gadget.Name, gadget.Output)
		if err != nil {
			log.Println(err)
		}
		g[i] = gadget
	}
	return g
}

func (a *App) Start(stop <-chan bool) {
	a.gadgets = append(a.gadgets, &Runner{})
	s := &Sockets{
		masterHost: a.masterHost,
		pubPort: a.pubPort,
		subPort: a.subPort,
	}
	a.gadgets = append(a.gadgets, s)
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

