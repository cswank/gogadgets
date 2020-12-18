package main

import "gopkg.in/alecthomas/kingpin.v2"

var (
	app = kingpin.New("gadgets", "gadgets")
	brk = app.Command("broker", "broker")
)

func main() {

}
