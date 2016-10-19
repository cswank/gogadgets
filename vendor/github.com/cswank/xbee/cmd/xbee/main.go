package main

import (
	"fmt"
	"log"
	"os"

	"github.com/cswank/xbee"
	"go.bug.st/serial.v1"
)

func main() {

	args := os.Args[1:]
	if len(args) < 1 {
		log.Fatal("you must pass in the /dev/tty port")
	}

	mode := &serial.Mode{}
	port, err := serial.Open(args[0], mode)
	if err != nil {
		log.Fatal(err)
	}

	for {
		msg := xbee.ReadMessage(port)
		a, err := msg.GetAnalog()
		if err != nil {
			log.Println(err)
		} else {
			for k, v := range a {
				fmt.Printf("%s: %.2f\n", k, v)
			}
		}
	}
}
