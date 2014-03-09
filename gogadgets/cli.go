package main


import (
	"flag"
	"time"
	"bitbucket.org/cswank/gogadgets"
	"os"
)

var (
	command = flag.String("c", "", "send a Robot Command Language command")
	host = flag.String("h", "localhost", "gadgets host")
)

func main() {
	flag.Parse()
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
