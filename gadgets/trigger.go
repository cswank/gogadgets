package output
/*
import (
	"fmt"
	"log"
	"time"
	"strings"
	"strconv"
	"bitbucket.com/cswank/gogadgets"
	"bitbucket.com/cswank/gogadgets/devices"
)



type Trigger struct {
	location string
	name string
	units string
	uid string
	operator string
	output devices.OutputDevice
	onCommand string
	offCommand string
	triggerType string
	compare Comparitor
	shutdown bool
	status bool	
	in <-chan gogadgets.Message
	out chan<- gogadgets.Message
	timerIn chan bool
	timerOut chan bool
}

func (t *Trigger) Start(out chan<- gogadgets.Message, in <-chan gogadgets.Message) {
	t.uid = fmt.Sprintf("%s %s trigger", t.location, t.name)
	t.in = in
	t.out = out
	t.timerIn = make(chan bool)
	t.timerOut = make(chan bool)
	for !t.shutdown {
		select {
		case msg := <-t.in:
			t.readMessage(&msg)
		case <-t.timerIn:
			t.output.Off()
			t.sendStatus()
		}
	}
}

*/
