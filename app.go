package main

import (
	//"log"
	"fmt"
	"bitbucket.com/cswank/gogadgets/models"
	//"bitbucket.com/cswank/gogadgets/gadgets"
)

type App struct {
	gadgets []models.Gadget
	channels []chan models.Message
}

func (a *App) Start(stop <-chan bool) {
	in := make(chan models.Message)
	a.channels = make([]chan models.Message, len(a.gadgets))
	for _, gadget := range a.gadgets {
		out := make(chan models.Message)
		go gadget.Start(out, in)
		a.channels = append(a.channels, out)
	}
	keepRunning := true
	fmt.Println("loop")
	for keepRunning {
		select {
		case msg := <-in:
			//fmt.Println("msg", msg)
			for _, channel := range a.channels {
				channel<- msg
			}
		case keepRunning = <-stop:
			stopMessage := models.Message{
				Type: models.COMMAND,
				Body: "shutdown",
			}
			for _, channel := range a.channels {
				channel<- stopMessage
			}
			for _, _  = range a.channels {
				<-in
			}
		}
	}
}

// func GadgetsFactory(configs []models.Config) []models.Gadget {
// 	g := make([]models.Gadget, len(configs))
// 	for i, config := range configs {
// 		gadget, err := gadgets.NewGadget(config)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 		g[i] = gadget
// 	}
// 	return g
// }

func main() {
	
}

