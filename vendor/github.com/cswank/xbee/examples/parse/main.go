package main

import (
	"fmt"
	"log"

	"github.com/cswank/xbee"
)

func main() {
	data := []byte{0x92, 0x00, 0x13, 0xA2, 0x00, 0x40, 0x4C, 0x0E, 0xBE, 0x61, 0x59, 0x01, 0x01, 0x00, 0x18, 0x03, 0x00, 0x10, 0x02, 0x2F, 0x01, 0xFE, 0x49}
	x, err := xbee.NewMessage(data)
	if err != nil {
		log.Fatal(err)
	}

	a, err := x.GetAnalog()
	if err != nil {
		log.Fatal(err)
	}

	d, err := x.GetDigital()
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range a {
		fmt.Println(k, v)
	}

	for k, v := range d {
		fmt.Println(k, v)
	}

}
