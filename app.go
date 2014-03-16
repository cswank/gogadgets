package gogadgets

import (
	"log"
	"io/ioutil"
	"encoding/json"
)

//App holds all the gadgets and handles passing Messages
//to them, and receiving messages from them.  It is the
//central part of Gadgets system.
type App struct {
	Gadgets    []GoGadget
	Host string
	PubPort    int
	SubPort    int
	channels   map[string]chan Message
	queue      *Queue
}	

//NewApp creates a new Gadgets system.  The cfg argument can be a
//path to a json file or a Config object itself.
func NewApp(cfg interface{}) *App {
	config := getConfig(cfg)
	if config.PubPort == 0 {
		config.SubPort = 6111
		config.PubPort = 6112
	}
	if config.Host == "" {
		config.Host = "localhost"
	}
	gadgets := GetGadgets(config.Gadgets)
	return &App{
		Host: config.Host,
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
	x := make(chan Message)
	a.GoStart(x)
}

//Useful for tests of other libraries that use Gogadgets
func (a *App) GoStart(input <-chan Message) {
	a.Gadgets = append(a.Gadgets, &Runner{})
	var sockets *Sockets
	sockets = &Sockets{
		host:    a.Host,
		pubPort: a.PubPort,
		subPort: a.SubPort,
	}
	a.Gadgets = append(a.Gadgets, sockets)
	in := make(chan Message)
	collect := make(chan Message)
	a.channels = make(map[string]chan Message)
	for _, gadget := range a.Gadgets {
		out := make(chan Message)
		a.channels[gadget.GetUID()] = out
		go gadget.Start(out, collect)
	}
	a.queue = NewQueue()
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

//Collects each message that is sent by the parts of the
//system.
func (a *App) collectMessages(in <-chan Message) {
	for {
		msg := <-in
		a.queue.Push(&msg)
	}
}

//After a message is collected by collectMessage, it is
//then sent back to the rest of the system.  This can 
//be improved.
func (a *App) dispenseMessages(out chan<- Message) {
	for {
		a.queue.cond.L.Lock()
		for a.queue.Len() == 0 {
			a.queue.Wait()
		}
		msg := a.queue.Get()
		out <- *msg
		a.queue.cond.L.Unlock()
	}
}

//This is where a new m
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


func getConfig(config interface{}) *Config {
	var c *Config
	switch v := config.(type) {
	case string:
		c = getConfigFromFile(v)
	case *Config:
		c = v
	default:
		panic("invalid config")
	}
	return c
}

func getConfigFromFile(configPath string) *Config {
	c := &Config{}
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(b, c)
	if err != nil {
		panic(err)
	}
	return c
}
