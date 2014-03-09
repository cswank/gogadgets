package gogadgets

import (
	"log"
	"time"
	"io/ioutil"
	"encoding/json"
)

//App holds all the gadgets and handles passing Messages
//to them, and receiving messages from them.  It is the
//central part of Gadgets system.
type App struct {
	Gadgets    []GoGadget
	MasterHost string
	PubPort    int
	SubPort    int
	channels   map[string]chan Message
	queue      *Queue
}

func NewApp(configPath string) *App {
	config := &Config{}
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, config)
	if err != nil {
		panic(err)
	}
	if config.PubPort == 0 {
		config.SubPort = 6111
		config.PubPort = 6112
	}
	if config.MasterHost == "" {
		config.MasterHost = "localhost"
	}
	gadgets := GetGadgets(config.Gadgets)
	return &App{
		MasterHost: config.MasterHost,
		PubPort:    config.PubPort,
		SubPort:    config.SubPort,
		Gadgets:    gadgets,
	}
}

//This is a factory fuction that reads a GadgtConfig
//and creates all the Gadgets that are defined in it.
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

//The main entry point for a Gadget system.  It takes
//a chan in case the system is started as a goroutine,
//but it can just be called directly.
func (a *App) Start() {
	a.Gadgets = append(a.Gadgets, &Runner{})
	sockets := &Sockets{
		host:    a.MasterHost,
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
		msg := <-in
		a.sendMessage(msg)
		if msg.Type == "command" && msg.Body == "shutdown" {
			keepRunning = false
			time.Sleep(100 * time.Millisecond)
		}
	}
}

//Collects each message that is sent by the parts of the
//system.
func (a *App) collectMessages(in <-chan Message) {
	for {
		msg := <-in
		a.queue.Push(&msg)
	}
}

//After a message is collected by collectMessage, it is
//then sent back out to each Gadget (except the Gadget
//that sent the message).
func (a *App) dispenseMessages(out chan<- Message) {
	for {
		if a.queue.Len() == 0 {
			time.Sleep(100 * time.Millisecond)
		} else {
			msg := a.queue.Get()
			out <- *msg
		}
	}
}

func (a *App) sendMessage(msg Message) {
	if msg.Target == "" {
		for uid, channel := range a.channels {
			if uid != msg.Sender {
				channel <- msg
			}
		}
	} else {
		channel, ok := a.channels[msg.Target]
		if ok {
			channel <- msg
		}
	}
}

//Some systems might have a few GoGadgets that are not
//built into the system (and hense can't be defined in
//the config file).  This is a way to apss in an instance
//of a gadget that is not part of the GoGadgets system.
func (a *App) AddGadget(gadget GoGadget) {
	a.Gadgets = append(a.Gadgets, gadget)
}
