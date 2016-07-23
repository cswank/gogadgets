package gogadgets

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

var (
	lg Logger
)

//App holds all the gadgets and handles passing Messages
//to them, and receiving Messages from them.  It is the
//central part of Gadgets system.
type App struct {
	gadgets []Gadgeter
	master  string
	host    string
	port    int
}

//NewApp creates a new Gadgets system.  The cfg argument can be a
//path to a json file or a Config object itself.
func NewApp(cfg interface{}, gadgets ...Gadgeter) *App {
	config := GetConfig(cfg)
	if config.Logger != nil {
		lg = config.Logger
	} else {
		lg = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	}
	if config.Port == 0 {
		config.Port = 6111
	}

	a := &App{
		master: config.Master,
		host:   config.Host,
		port:   config.Port,
	}
	a.GetGadgets(config.Gadgets)
	a.gadgets = append(a.gadgets, gadgets...)
	return a
}

//GetGadgets is a factory fuction that reads a GadgtConfig
//and creates all the Gadgets that are defined in it.
func (a *App) GetGadgets(configs []GadgetConfig) {
	a.gadgets = make([]Gadgeter, len(configs))
	for i, config := range configs {
		gadget, err := NewGadget(&config)
		if err != nil {
			log.Fatal(err)
		}
		a.gadgets[i] = gadget
	}
	a.gadgets = append(a.gadgets, &MethodRunner{})
	srv := NewServer(a.host, a.master, a.port, lg)
	a.gadgets = append(a.gadgets, srv)

}

//Start is the main entry point for a Gadget system.  It takes
//a chan in case the system is started as a goroutine,
//but it can just be called directly.
func (a *App) Start() {
	x := make(chan Message)
	a.GoStart(x)
}

// GoStart enables a gadgets system to be started
// by either a test suite that needs it to run
// as a goroutine or a client app that starts
// gogadget systems upon a command from a central
// web app.
func (a *App) GoStart(input <-chan Message) {

	collect := make(chan Message)
	channels := make(map[string]chan Message)
	for _, gadget := range a.gadgets {
		out := make(chan Message)
		channels[gadget.GetUID()] = out
		go gadget.Start(out, collect)
	}
	lg.Println("started gagdgets")
	b := NewBroker(channels, input, collect)
	b.Start()
}

func GetConfig(config interface{}) *Config {
	var c *Config
	switch v := config.(type) {
	case string:
		c = getConfigFromFile(v)
	case *Config:
		c = v
	default:
		lg.Fatal("invalid config")
	}
	return c
}

func getConfigFromFile(configPath string) *Config {
	c := &Config{}
	b, err := ioutil.ReadFile(configPath)
	if err != nil {
		lg.Fatal(err)
	}
	err = json.Unmarshal(b, c)
	if err != nil {
		lg.Fatal(err)
	}
	return c
}
