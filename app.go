package gogadgets

import (
	"encoding/json"
	"fmt"
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
	Gadgets []Gadgeter
	Master  string
	Host    string
	Port    int
}

//NewApp creates a new Gadgets system.  The cfg argument can be a
//path to a json file or a Config object itself.
func NewApp(cfg interface{}) *App {
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
		Master: config.Master,
		Host:   config.Host,
		Port:   config.Port,
	}
	a.GetGadgets(config.Gadgets)
	return a
}

//This is a factory fuction that reads a GadgtConfig
//and creates all the Gadgets that are defined in it.
func (a *App) GetGadgets(configs []GadgetConfig) {
	a.Gadgets = make([]Gadgeter, len(configs))
	for i, config := range configs {
		gadget, err := NewGadget(&config)
		if err != nil {
			a.Gadgets[i] = &Error{
				location: config.Location,
				name:     config.Name,
				error:    fmt.Errorf("couldn't initialize %s %s: %s\n", config.Location, config.Name, err),
			}
		} else {
			a.Gadgets[i] = gadget
		}
	}
	a.Gadgets = append(a.Gadgets, &MethodRunner{})
	srv := NewServer(a.Host, a.Master, a.Port, lg)
	a.Gadgets = append(a.Gadgets, srv)
}

//The main entry point for a Gadget system.  It takes
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
	for _, gadget := range a.Gadgets {
		out := make(chan Message)
		channels[gadget.GetUID()] = out
		go gadget.Start(out, collect)
	}
	lg.Println("started gagdgets")
	b := NewBroker(channels, input, collect)
	b.Start()
}

//Some systems might have a few GoGadgets that are not
//built into the system (and hence can't be defined in
//the config file).  This is a way to add in an instance
//of a gadget that is not part of the GoGadgets system.
func (a *App) AddGadget(gadget Gadgeter) {
	a.Gadgets = append(a.Gadgets, gadget)
}

func GetConfig(config interface{}) *Config {
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
