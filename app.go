package gogadgets

import (
	"log"
	"time"
)

type App struct {
	Gadgets []GoGadget
	MasterHost string
	PubPort int
	SubPort int
	channels map[string]chan Message
	queue *Queue
}

func NewApp(config *Config) *App {
	gadgets := GetGadgets(config.Gadgets)
	if  config.PubPort == 0 {
		config.SubPort = 6111
		config.PubPort = 6112
	}
	return &App{
		MasterHost: config.MasterHost,
		PubPort: config.PubPort,
		SubPort: config.SubPort,
		Gadgets: gadgets,
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

func (a *App) Start(input <-chan Message) {
	a.Gadgets = append(a.Gadgets, &Runner{})
	sockets := &Sockets{
		host: a.MasterHost,
		pubPort: a.PubPort,
		subPort: a.SubPort,
	}
	a.Gadgets = append(a.Gadgets, sockets)
	in := make(chan Message)
	collect := make(chan Message)
	a.channels = make(map[string]chan Message)
	a.queue = NewQueue()
	for _, gadget := range a.Gadgets {
		out := make(chan Message)
		a.channels[gadget.GetUID()] = out
		go gadget.Start(out, collect)
	}
	go a.collectMessages(collect)
	go a.dispenseMessages(in)
	keepRunning := true
	log.Println("started gagdgets")
	for keepRunning {
		select {
		case msg := <-in:
			a.sendMessage(msg)
		case msg := <-input:
			if msg.Type == "command" && msg.Body == "shutdown" {
				keepRunning = false
			}
			a.sendMessage(msg)
		}
	}
}

func (a *App) collectMessages(in <-chan Message) {
	for {
		msg := <-in
		a.queue.Push(&msg)
	}
}

func (a *App) dispenseMessages(out chan<- Message) {
	for {
		if a.queue.Len() == 0 {
			time.Sleep(100 * time.Millisecond)
		} else {
			msg := a.queue.Get()
			out<- *msg
		}
	}
}

func (a *App) sendMessage(msg Message) {
	if msg.Target == "" {
		for uid, channel := range a.channels {
			if uid != msg.Sender {
				channel<- msg
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
	a.Gadgets = append(a.Gadgets, gadget)
}
