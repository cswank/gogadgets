package main

import (
	"fmt"
	"bitbucket.com/cswank/gogadgets"
)

func main() {
	cfg := &gogadgets.Config{
		MasterHost: "localhost",
		Gadgets: []gogagets.GadgetConfig{
			gogadgets.GadgetConfig{
				Type: "switch",
				Location: "tank",
				Name: "volume",
				Pin: gogadgets.Pin{
					Port: "9",
					Pin: "16",
					Value: 5.0,
					Units: "liters",
				},
			},
			gogadgets.GadgetConfig{
				Type: "gpio",
				Location: "tank",
				Name: "water",
				Pin: gogadgets.Pin{
					Port: "9",
					Pin: "15",
				},
			},
			gogadgets.GadgetConfig{
				Type: "gpio",
				Location: "lab",
				Name: "led",
				Pin: gogadgets.Pin{
					Port: "9",
					Pin: "14",
				},
			},
			gogadgets.GadgetConfig{
				Type: "temperature",
				Location: "tank",
				Name: "temperature",
				Pin: gogadgets.Pin{
					OneWireId: "28-0000047ade8f",
					Units: "C",
				},
			},
		},
	}
	app := gogadgets.NewApp(cfg)
	stop := make(chan bool)
	app.Start(stop)
}
