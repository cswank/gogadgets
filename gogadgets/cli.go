package main


import (
	"flag"
	"time"
	"fmt"
	"bitbucket.org/cswank/gogadgets"
	"os"
)

var (
	host = flag.String("h", "localhost", "gadgets host (for the send command)")
	//config = flag.String("c", "", "config file for running a gadgets system")
)

func main() {
	flag.Parse()
	sub := flag.Arg(0)
	if sub == "run" {
		runGadgets()
	} else if sub == "send" {
		sendCommand()
	} else if sub == "listen" {
		listen()
	}
}

func runGadgets() {
	cfg := flag.Arg(1)
	a := gogadgets.NewApp(cfg)
	a.Start()
}

func sendCommand() {	
	s, err := gogadgets.NewClientSockets(*host)
	defer s.Close()
	if err != nil {
		panic(err)
	}
	time.Sleep(10 * time.Millisecond)
	fmt.Println(flag.Arg(1), *host)
	s.Send(flag.Arg(1))
	time.Sleep(10 * time.Millisecond)
	os.Exit(0)
}

//Waits for a zmq message that contains a gogadgets
//config.  When one is recieved it is parsed and a
//a gogadgts system is started.
func listen() {
	
}
