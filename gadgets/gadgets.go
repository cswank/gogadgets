package output

import (
	"log"
	"time"
	"fmt"
	"strings"
	"strconv"
	"bitbucket.com/cswank/gogadgets"
	"bitbucket.com/cswank/gogadgets/devices"
)


var (
	units = map[string]string{
		"liters": "volume",
		"gallons": "volume",
		"liter": "volume",
		"gallon": "volume",
		"c": "temperature",
		"f": "temperature",
		"celcius": "temperature",
		"fahrenheit": "temperature",
		"seconds": "time",
		"minutes": "time",
		"hours": "time",
		"second": "time",
		"minute": "time",
		"hour": "time",
	}
)

type Comparitor func(msg *gogadgets.Message) bool

type Config struct {
	Location string
	Name string
	Pin *devices.Pin
}

type Gadget struct {
	gogadgets.Gadget
	location string
	name string
	output devices.OutputDevice
	onCommand string
	offCommand string
	uid string
	status bool
	compare Comparitor
	shutdown bool
	units string
	operator string
	out chan<- gogadgets.Message
	timerIn chan bool
	timerOut chan bool
}

func NewOutputGadget(config *Config) (*Gadget, error) {
	dev, err := devices.NewOutputDevice(config.Pin)
	if err == nil {
		return &Gadget{
			location: config.Location,
			name: config.Name,
			onCommand: fmt.Sprintf("turn on %s %s", config.Location, config.Name),
			offCommand: fmt.Sprintf("turn off %s %s", config.Location, config.Name),
			output: dev,
			uid: fmt.Sprintf("%s %s", config.Location, config.Name),
		}, err
	}
	return nil, err
}

func (g *Gadget) isMyCommand(msg *gogadgets.Message) bool {
	return msg.Type == gogadgets.COMMAND && 
		(strings.Index(msg.Body, g.onCommand) == 0 ||
		strings.Index(msg.Body, g.offCommand) == 0 ||
		msg.Body == "shutdown")
}

func (g *Gadget) Start(in <-chan gogadgets.Message, out chan<- gogadgets.Message) {
	g.out = out
	g.timerIn = make(chan bool)
	g.timerOut = make(chan bool)
	for !g.shutdown {
		select {
		case msg := <-in:
			g.readMessage(&msg)			
		case <-g.timerOut:
			g.off()
		}
	}
}

func (g *Gadget) off() {
	g.status = false
	g.output.Off()
	g.compare = nil
	g.sendStatus()
}

func (g *Gadget) on(val *gogadgets.Value) {
	g.output.On(val)
	if !g.status {
		g.status = true
		g.sendStatus()
	}
}

func (g *Gadget) readMessage(msg *gogadgets.Message) {
	if msg.Type == gogadgets.COMMAND && g.isMyCommand(msg) {
		g.readCommand(msg)
	} else if g.status && msg.Type == gogadgets.STATUS {
		g.readStatus(msg)
	}
}

func (g *Gadget) readStatus(msg *gogadgets.Message) {
	if g.status && g.compare != nil && g.compare(msg) {
		g.off()
	} else if g.status && msg.Location == g.location {
		g.output.Update(msg)
	}
}

func (g *Gadget) readCommand(msg *gogadgets.Message) {
	if msg.Body == "shutdown" {
		g.shutdown = true
		g.off()
	} else if strings.Index(msg.Body, g.onCommand) == 0 {
		g.readOnCommand(msg)
	} else if strings.Index(msg.Body, g.offCommand) == 0 {
		g.readOffCommand(msg)
	}
}

func (g *Gadget) readOnCommand(msg *gogadgets.Message) {
	var val *gogadgets.Value
	if len(strings.Trim(msg.Body, " ")) > len(g.onCommand) {
		val = g.readOnArguments(msg.Body)
	} else {
		g.compare = nil
		
	}
	g.on(val)
}

func (g *Gadget) readOnArguments(cmd string) *gogadgets.Value {
	var val *gogadgets.Value
	value, unit, err := g.getValue(cmd)
	if err != nil {
		log.Println("could not parse", cmd)
	}
	gadget, ok := units[unit]
	if ok {
		if gadget == "time" {
			go g.startTimer(value, unit, g.timerIn, g.timerOut)
		} else if gadget == "volume" || gadget == "temperature" {
			g.setCompare(value, unit, gadget)
			val = &gogadgets.Value{
				Value: value,
				Units: unit,
			}
		}
	}
	return val
}

func (g *Gadget) setCompare(value float64, unit string, gadget string) {
	if g.operator == "<=" {
		g.compare = func(msg *gogadgets.Message) bool {
			val, ok := msg.Value.Value.(float64)
			return msg.Location == g.location &&
				ok &&
				msg.Name == gadget &&
				val <= value
		}
	} else if g.operator == ">=" {
		g.compare = func(msg *gogadgets.Message) bool {
			val, ok := msg.Value.Value.(float64)
			return msg.Location == g.location &&
				ok &&
				msg.Name == gadget &&
				val >= value
		}
	}
}

func (g *Gadget) getValue(cmd string) (float64, string, error) {
	cmd = g.stripCommand(cmd)
	value, unit, err := g.splitCommand(cmd)
	var v float64
	if err == nil {
		v, err = strconv.ParseFloat(value, 64)
	}
	return v, unit, err
}

func (g *Gadget) startTimer(value float64, unit string, in <-chan bool, out chan<- bool) {
	d := time.Duration(value * float64(time.Second))
	keepGoing := true
	for  keepGoing{
		select {
		case <-in:
			keepGoing = false
		case <-time.After(d):
			keepGoing = false
			out<- true
		}
	}
}

func (g *Gadget) splitCommand(cmd string) (string, string, error) {
	parts := strings.Split(cmd, " ")
	return parts[0], parts[1], nil
}

func (g *Gadget) stripCommand(cmd string) string {
	cmd = strings.Trim(cmd, " ")
	cmd = strings.TrimPrefix(cmd, g.onCommand)
	cmd = strings.TrimPrefix(cmd, " for ")
	return strings.TrimPrefix(cmd, " to ")
}

func (g *Gadget) readOffCommand(msg *gogadgets.Message) {
	if g.status {
		g.off()
	}
}

func (g *Gadget) sendStatus() {
	msg := gogadgets.Message{
		Sender: g.uid,
		Type: gogadgets.STATUS,
		Location: g.location,
		Name: g.name,
		Value: gogadgets.Value{
			Units: g.units,
			Value:g.status,
		},
	}
	g.out<- msg
}

