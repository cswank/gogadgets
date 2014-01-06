package gogadgets

import (
	"time"
	"fmt"
	"errors"
	"strings"
	"strconv"
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

type Comparitor func(msg *Message) bool

type Gadget struct {
	GoGadget
	Location string `json:"location"`
	Name string `json:"name"`
	Output OutputDevice `json:"-"`
	Input InputDevice `json:"-"`
	Direction string `json:"direction"`
	OnCommand string `json:"on"`
	OffCommand string `json:"off"`
	UID string `json:"uid"`
	status bool
	compare Comparitor
	shutdown bool
	units string
	Operator string
	out chan<- Message
	devIn chan Message
	timerIn chan bool
	timerOut chan bool
}

func NewGadget(config *GadgetConfig) (*Gadget, error) {
	t := config.Pin.Type
	if t == "heater" || t == "cooler" || t == "gpio" {
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

func NewInputGadget(config *GadgetConfig) (gadget *Gadget, err error) {
	dev, err := NewInputDevice(&config.Pin)
	if err == nil {
		gadget = &Gadget{
			Location: config.Location,
			Name: config.Name,
			Input: dev,
			Direction: "input",
			OnCommand: "n/a",
			OffCommand: "n/a",
			UID: fmt.Sprintf("%s %s", config.Location, config.Name),
		}
	}
	return gadget, err
}

func NewOutputGadget(config *GadgetConfig) (gadget *Gadget, err error) {
	dev, err := NewOutputDevice(&config.Pin)
	if config.OnCommand == "" {
		config.OnCommand = fmt.Sprintf("turn on %s %s", config.Location, config.Name)
	}
	if config.OffCommand == "" {
		config.OffCommand = fmt.Sprintf("turn off %s %s", config.Location, config.Name)
	}
	if err == nil {
		gadget = &Gadget{
			Location: config.Location,
			Name: config.Name,
			Direction: "output",
			OnCommand: config.OnCommand,
			OffCommand: config.OffCommand,
			Output: dev,
			UID: fmt.Sprintf("%s %s", config.Location, config.Name),
		}
	} else {
		panic(err)
	}
	return gadget, err
}

func (g *Gadget) isMyCommand(msg *Message) bool {
	return msg.Type == COMMAND && 
		(strings.Index(msg.Body, g.OnCommand) == 0 ||
		strings.Index(msg.Body, g.OffCommand) == 0 ||
		msg.Body == "update" ||
		msg.Body == "shutdown")
}

func (g *Gadget) Start(in <-chan Message, out chan<- Message) {
	g.out = out
	g.timerIn = make(chan bool)
	g.timerOut = make(chan bool)
	if g.Output != nil {
		g.off()
		g.doOutputLoop(in)
	} else if g.Input != nil {
		g.doInputLoop(in)
	}
}

func (g *Gadget) doInputLoop(in <-chan Message) {
	devOut := make(chan Value)
	g.devIn = make(chan Message)
	go g.Input.Start(g.devIn, devOut)
	for !g.shutdown {
		select {
		case msg := <-in:
			g.readMessage(&msg)
		case val := <-devOut:
			g.out<- Message{
				Sender: g.UID,
				Type: "update",
				Location: g.Location,
				Name: g.Name,
				Value: val,
			}
		}
	}
}

func (g *Gadget) doOutputLoop(in <-chan Message) {
	for !g.shutdown {
		select {
		case msg := <-in:
			g.readMessage(&msg)			
		case <-g.timerOut:
			g.off()
		}
	}
}

func (g *Gadget) on(val *Value) {
	g.Output.On(val)
	if !g.status {
		g.status = true
		go g.sendUpdate()
	}
}

func (g *Gadget) off() {
	g.status = false
	g.Output.Off()
	g.compare = nil
	go g.sendUpdate()
}


func (g *Gadget) readMessage(msg *Message) {
	if g.Input != nil {
		fmt.Println("sending msg to input dev")
		g.devIn<- *msg
	}
	if msg.Type == COMMAND && g.isMyCommand(msg) {
		g.readCommand(msg)
	} else if g.status && msg.Type == UPDATE {
		g.readUpdate(msg)
	}
}

func (g *Gadget) readUpdate(msg *Message) {
	if g.status && g.compare != nil && g.compare(msg) {
		g.off()
	} else if g.status && msg.Location == g.Location {
		g.Output.Update(msg)
	}
}

func (g *Gadget) readCommand(msg *Message) {
	if msg.Body == "shutdown" {
		g.shutdown = true
		g.off()
	} else if msg.Body == "update" {
		go g.sendUpdate()
	} else if strings.Index(msg.Body, g.OnCommand) == 0 {
		g.readOnCommand(msg)
	} else if strings.Index(msg.Body, g.OffCommand) == 0 {
		g.readOffCommand(msg)
	}
}

func (g *Gadget) readOnCommand(msg *Message) {
	var val *Value
	if len(strings.Trim(msg.Body, " ")) > len(g.OnCommand) {
		val, err := g.readOnArguments(msg.Body)
		if err == nil {
			g.on(val)
		}
	} else {
		g.compare = nil
		g.on(val)
	}
}

func (g *Gadget) readOnArguments(cmd string) (*Value, error) {
	var val *Value
	value, unit, err := ParseCommand(cmd)
	if err != nil {
		return val, errors.New(fmt.Sprintf("could not parse %s", cmd))
	}
	gadget, ok := units[unit]
	if ok {
		if gadget == "time" {
			go g.startTimer(value, unit, g.timerIn, g.timerOut)
		} else if gadget == "volume" || gadget == "temperature" {
			g.setCompare(value, unit, gadget)
			val = &Value{
				Value: value,
				Units: unit,
			}
		}
	}
	return val, nil
}

func (g *Gadget) setCompare(value float64, unit string, gadget string) {
	if g.Operator == "<=" {
		g.compare = func(msg *Message) bool {
			val, ok := msg.Value.Value.(float64)
			return msg.Location == g.Location &&
				ok &&
				msg.Name == gadget &&
				val <= value
		}
	} else if g.Operator == ">=" {
		g.compare = func(msg *Message) bool {
			val, ok := msg.Value.Value.(float64)
			return msg.Location == g.Location &&
				ok &&
				msg.Name == gadget &&
				val >= value
		}
	}
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

func (g *Gadget) readOffCommand(msg *Message) {
	if g.status {
		g.off()
	}
}

func (g *Gadget) GetUID() string {
	if g.UID == "" {
		g.UID = fmt.Sprintf("%s %s", g.Location, g.Name)
	}
	return g.UID
}

func (g *Gadget) sendUpdate() {
	var value *Value
	if g.Input != nil {
		value = g.Input.GetValue()
	} else {
		value = &Value{
			Units: g.units,
			Value: g.status,
		}
	}
	msg := Message{
		Sender: g.UID,
		Type: UPDATE,
		Location: g.Location,
		Name: g.Name,
		Value: *value,
		Info: Info{
			Direction: g.Direction,
			On: g.OnCommand,
			Off: g.OffCommand,
		},
	}
	g.out<- msg
}

func ParseCommand(cmd string) (float64, string, error) {
	cmd = stripCommand(cmd)
	value, unit, err := splitCommand(cmd)
	var v float64
	if err == nil {
		v, err = strconv.ParseFloat(value, 64)
	}
	return v, unit, err
}

func splitCommand(cmd string) (string, string, error) {
	parts := strings.Split(cmd, " ")
	return parts[0], parts[1], nil
}

func stripCommand(cmd string) string {
	cmd = strings.Trim(cmd, " ")
	i := strings.Index(cmd, " for ")
	if i != -1 {
		return cmd[i + 5:]
	}
	i = strings.Index(cmd, " to ")
	fmt.Println("index is", i)
	if i != -1 {
		return cmd[i + 4:]
	}
	return ""
}
