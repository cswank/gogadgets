package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/cswank/gogadgets"
	"github.com/cswank/gogadgets/utils"
)

const (
	defaultDir    = "~/.gadgets"
	defaultConfig = "/Users/Cswank/.gadgets/config.json"
)

var (
	host   = flag.String("h", "localhost", "Name of Host")
	config = flag.String("g", "", "Path to a Gadgets config file")
	cmd    = flag.String("c", "", "a Robot Command Language string")
	status = flag.Bool("s", false, "get the status of a gadgets system")
)

func main() {
	flag.Parse()
	if len(*cmd) > 0 {
		sendCommand()
	} else if *status {
		getStatus()
	} else {
		runGadgets()
	}
}

func runGadgets() {
	cfg := getConfig()
	if cfg == "" {
		listen()
	} else {
		a := gogadgets.NewApp(cfg)
		a.Start()
	}
}

func getConfig() string {
	cfg := *config
	if cfg != "" {
		return cfg
	}
	if utils.FileExists(defaultConfig) {
		return defaultConfig
	}
	return ""
}

func getStatus() {
	cfg := gogadgets.SocketsConfig{
		Host:    *host,
		SubPort: 6111,
		PubPort: 6112,
		Master:  false,
	}
	s, err := gogadgets.NewClientSockets(cfg)
	defer s.Close()
	if err != nil {
		panic(err)
	}
	time.Sleep(100 * time.Millisecond)
	status, err := s.SendStatusRequest()
	time.Sleep(100 * time.Millisecond)
	if err == nil {
		fmt.Println("status", status, err)
	} else {
		fmt.Println(err)
	}
	os.Exit(0)
}

func sendCommand() {
	cfg := gogadgets.SocketsConfig{
		Host:    *host,
		SubPort: 6111,
		PubPort: 6112,
		Master:  false,
	}
	s, err := gogadgets.NewClientSockets(cfg)
	defer s.Close()
	if err != nil {
		panic(err)
	}
	time.Sleep(100 * time.Millisecond)
	fmt.Println(*cmd, "host", *host)
	s.Send(*cmd)
	time.Sleep(100 * time.Millisecond)
	os.Exit(0)

}

//Waits for a zmq message that contains a gogadgets
//config.  When one is recieved it is written to the
//default config path and a a gogadgts system is started.
func listen() {
	cfg := gogadgets.SocketsConfig{
		Host:    *host,
		SubPort: 6111,
		PubPort: 6112,
		Master:  false,
	}
	s, err := gogadgets.NewSockets(cfg)
	if err != nil {
		panic(err)
	}
	time.Sleep(100 * time.Millisecond)
	log.Println("listening for new gadgets")
	msg := s.Recv()
	d, err := json.Marshal(&msg.Config)
	if err != nil {
		panic(err)
	}
	os.Mkdir(defaultDir, 0644)
	err = ioutil.WriteFile(defaultConfig, d, 0644)
	if err != nil {
		panic(err)
	}
	s.Close()
	time.Sleep(100 * time.Millisecond)
	runGadgets()
}
