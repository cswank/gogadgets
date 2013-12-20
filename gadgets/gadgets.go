package gadgets

import (
	"log"
	"time"
	"fmt"
	"errors"
	"strings"
	"strconv"
	"bitbucket.com/cswank/gogadgets/models"
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

type Comparitor func(msg *models.Message) bool

type Gadget struct {
	models.Gadget
	Location string
	Name string
	Output devices.OutputDevice
	Input devices.InputDevice
	OnCommand string
	OffCommand string
	UID string
	status bool
	compare Comparitor
	shutdown bool
	units string
	Operator string
	out chan<- models.Message
	timerIn chan bool
	timerOut chan bool
}

func NewGadget(config *models.Config) (*Gadget, error) {
	t := config.Pin.Type
	if t == "heater" || t == "gpio" {
		return NewOutputGadget(config)
	} else if t == "thermometer" || t == "switch" {
		return NewInputGadget(config)
	}
	err := errors.New(
		fmt.Sprintf(
			"couldn't build a gadget based on config: %s %s",
			config.Location,
			config.Name))
	return nil, err
}

func NewInputGadget(config *models.Config) (gadget *Gadget, err error) {
	dev, err := devices.NewInputDevice(&config.Pin)
	if err == nil {
		gadget = &Gadget{
			Location: config.Location,
			Name: config.Name,
			Input: dev,
			UID: fmt.Sprintf("%s %s", config.Location, config.Name),
		}
	}
	return gadget, err
}

func NewOutputGadget(config *models.Config) (gadget *Gadget, err error) {
	dev, err := devices.NewOutputDevice(&config.Pin)
	if err == nil {
		gadget = &Gadget{
			Location: config.Location,
			Name: config.Name,
			OnCommand: fmt.Sprintf("turn on %s %s", config.Location, config.Name),
			OffCommand: fmt.Sprintf("turn off %s %s", config.Location, config.Name),
			Output: dev,
			UID: fmt.Sprintf("%s %s", config.Location, config.Name),
		}
	}
	return gadget, err
}

func (g *Gadget) isMyCommand(msg *models.Message) bool {
	return msg.Type == models.COMMAND && 
		(strings.Index(msg.Body, g.OnCommand) == 0 ||
		strings.Index(msg.Body, g.OffCommand) == 0 ||
		msg.Body == "shutdown")
}

func (g *Gadget) Start(in <-chan models.Message, out chan<- models.Message) {
	g.out = out
	g.timerIn = make(chan bool)
	g.timerOut = make(chan bool)
	if g.Output != nil {
		g.doOutputLoop(in)
	} else if g.Input != nil {
		g.doInputLoop(in)
	}
}

func (g *Gadget) doInputLoop(in <-chan models.Message) {
	devOut := make(chan models.Value)
	stop := make(chan bool)
	go g.Input.Start(stop, devOut)
	for !g.shutdown {
		select {
		case msg := <-in:
			g.readMessage(&msg)
		case val := <-devOut:
			g.out<- models.Message{
				Location: g.Location,
				Name: g.Name,
				Value: val,
			}
		}
	}
}

func (g *Gadget) doOutputLoop(in <-chan models.Message) {
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
	g.Output.Off()
	g.compare = nil
	g.sendStatus()
}

func (g *Gadget) on(val *models.Value) {
	g.Output.On(val)
	if !g.status {
		g.status = true
		g.sendStatus()
	}
}

func (g *Gadget) readMessage(msg *models.Message) {
	if msg.Type == models.COMMAND && g.isMyCommand(msg) {
		g.readCommand(msg)
	} else if g.status && msg.Type == models.STATUS {
		g.readStatus(msg)
	}
}

func (g *Gadget) readStatus(msg *models.Message) {
	if g.status && g.compare != nil && g.compare(msg) {
		g.off()
	} else if g.status && msg.Location == g.Location {
		g.Output.Update(msg)
	}
}

func (g *Gadget) readCommand(msg *models.Message) {
	if msg.Body == "shutdown" {
		g.shutdown = true
		g.off()
	} else if strings.Index(msg.Body, g.OnCommand) == 0 {
		g.readOnCommand(msg)
	} else if strings.Index(msg.Body, g.OffCommand) == 0 {
		g.readOffCommand(msg)
	}
}

func (g *Gadget) readOnCommand(msg *models.Message) {
	var val *models.Value
	if len(strings.Trim(msg.Body, " ")) > len(g.OnCommand) {
		val = g.readOnArguments(msg.Body)
	} else {
		g.compare = nil
		
	}
	g.on(val)
}

func (g *Gadget) readOnArguments(cmd string) *models.Value {
	var val *models.Value
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
			val = &models.Value{
				Value: value,
				Units: unit,
			}
		}
	}
	return val
}

func (g *Gadget) setCompare(value float64, unit string, gadget string) {
	if g.Operator == "<=" {
		g.compare = func(msg *models.Message) bool {
			val, ok := msg.Value.Value.(float64)
			return msg.Location == g.Location &&
				ok &&
				msg.Name == gadget &&
				val <= value
		}
	} else if g.Operator == ">=" {
		g.compare = func(msg *models.Message) bool {
			val, ok := msg.Value.Value.(float64)
			return msg.Location == g.Location &&
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
	cmd = strings.TrimPrefix(cmd, g.OnCommand)
	cmd = strings.TrimPrefix(cmd, " for ")
	return strings.TrimPrefix(cmd, " to ")
}

func (g *Gadget) readOffCommand(msg *models.Message) {
	if g.status {
		g.off()
	}
}

func (g *Gadget) sendStatus() {
	if g.UID == "" {
		g.UID = fmt.Sprintf("%s %s", g.Location, g.Name)
	}
	msg := models.Message{
		Sender: g.UID,
		Type: models.STATUS,
		Location: g.Location,
		Name: g.Name,
		Value: models.Value{
			Units: g.units,
			Value: g.status,
		},
	}
	g.out<- msg
}

