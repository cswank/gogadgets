package main


import (
	"flag"
	"time"
	"bitbucket.org/cswank/gogadgets"
	"os"
)

var (
	command = flag.String("c", "", "send a Robot Command Language command")
	cfg = flag.String("config", "", "Path to the config json file")
	host = flag.String("h", "localhost", "gadgets host")
)

func main() {
	flag.Parse()
	if len(*cfg ) > 0 {
		runGadgets()
	} else if len(*command) > 0 {
		sendCommand()
	}
}

func runGadgets() {
	a := gogadgets.NewApp(cfg)
	a.Start()
}

func sendCommand() {	
	s, err := gogadgets.NewClientSockets(*host)
	defer s.Close()
	if err != nil {
		panic(err)
	}
	time.Sleep(50 * time.Millisecond)
	s.Send(*command)
	time.Sleep(50 * time.Millisecond)
	os.Exit(0)
}
