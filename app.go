package gogadgets

import (
	"encoding/json"
	"log"
	"os"
)

// App holds all the gadgets and handles passing Messages
// to them, and receiving Messages from them.  It is the
// central part of Gadgets system.
type App struct {
	gadgets []Gadgeter
}

// New creates a new Gadgets system.  The cfg argument can be a
// path to a json file or a Config object itself.
func New(cfg interface{}, gadgets ...Gadgeter) *App {
	config := GetConfig(cfg)
	return &App{
		gadgets: config.CreateGadgets(gadgets...),
	}
}

// Start is the main entry point for a Gadget system.  It takes
// a chan in case the system is started as a goroutine,
// but it can just be called directly.
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
	log.Println("started gagdgets")
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
		log.Fatal("invalid config")
	}

	if c.Port == 0 {
		c.Port = 6111
	}
	return c
}

func getConfigFromFile(pth string) *Config {
	c := &Config{}
	f, err := os.Open(pth)
	if err != nil {
		log.Fatalf("unable to open config path %s: %s", pth, err)
	}

	defer f.Close()

	err = json.NewDecoder(f).Decode(c)
	if err != nil {
		log.Fatalf("unable to parse json from config path %s: %s", pth, err)
	}
	return c
}
