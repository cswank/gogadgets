package main

import (
	"bitbucket.com/cswank/gogadgets"
)

func main() {
	cfg := &gogadgets.Config{
		MasterHost: "localhost",
		Gadgets: []gogadgets.GadgetConfig{
			gogadgets.GadgetConfig{
				Location: "tank",
				Name: "volume",
				Pin: gogadgets.Pin{
					Type: "switch",
					Port: "9",
					Pin: "16",
					Value: 5.0,
					Units: "liters",
				},
			},
			gogadgets.GadgetConfig{
				Location: "tank",
				Name: "water",
				Pin: gogadgets.Pin{
					Type: "gpio",
					Port: "9",
					Pin: "15",
				},
			},
			gogadgets.GadgetConfig{
				
				Location: "lab",
				Name: "led",
				Pin: gogadgets.Pin{
					Type: "gpio",
					Port: "9",
					Pin: "14",
				},
			},
			gogadgets.GadgetConfig{
				Location: "tank",
				Name: "temperature",
				Pin: gogadgets.Pin{
					Type: "thermometer",
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
