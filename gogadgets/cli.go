package main

import (
	"flag"
	"time"
	"fmt"
	"bitbucket.org/cswank/gogadgets"
	"os"
)

var (
	host = flag.String("h", "localhost", "Name of Host")
	config = flag.String("g", "", "Path to a Gadgets config file")
	cmd = flag.String("c", "", "a Robot Command Language string")
)

func main() {
	flag.Parse()
	if len(*config) > 0 {
		runGadgets()
	} else if len(*cmd) > 0 {
		sendCommand()
	} else {
		listen()
	}
}

func runGadgets() {
	a := gogadgets.NewApp(*config)
	a.Start()
}

func sendCommand() {	
	s, err := gogadgets.NewClientSockets(*host)
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
//config.  When one is recieved it is parsed and a
//a gogadgts system is started.
func listen() {
	s, err := gogadgets.NewSockets()
	if err != nil {
		panic(err)
	}
	time.Sleep(1 * time.Second)
	fmt.Println("waiting for message")
	msg := s.Recv()
	fmt.Println("got a msg", msg)
	s.Close()
}
