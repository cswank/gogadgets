package main



import (
	"github.com/droundy/goopt"
	"time"
	"fmt"
	"bitbucket.org/cswank/gogadgets"
	"os"
)

var (
	host = goopt.String([]string{"-h", "--host"}, "localhost", "Name of Host")
	config = goopt.String([]string{"-g", "--gadgets"}, "", "Path to a Gadgets config file")
	cmd = goopt.String([]string{"-c", "--cmd"}, "", "a Robot Command Language string")
)

func main() {
	goopt.Parse(nil)
	fmt.Println(len(*config))
	if len(*config) > 0 {
		runGadgets()
	} else if len(*cmd) > 0 {
		sendCommand()
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
	
}
