package gogadgets

import (
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

func Init(p PortFactory) {
	serialFactory = p
}

var (
	units = map[string]string{
		"liters":     "volume",
		"gallons":    "volume",
		"liter":      "volume",
		"gallon":     "volume",
		"c":          "temperature",
		"f":          "temperature",
		"C":          "temperature",
		"F":          "temperature",
		"celcius":    "temperature",
		"fahrenheit": "temperature",
		"seconds":    "time",
		"minutes":    "time",
		"hours":      "time",
		"second":     "time",
		"minute":     "time",
		"hour":       "time",
		"%":          "power",
	}
)

type Comparitor func(msg *Message) bool

type Gadgeter interface {
	GetUID() string
	GetDirection() string
	Start(in <-chan Message, out chan<- Message)
}

//Each part of a Gadgets system that controls a single
//piece of hardware (for example: a gpio pin) is represented
//by Gadget.  A Gadget must have either an InputDevice or
//an OutputDevice.  Gadget fulfills the GoGaget interface.
type Gadget struct {
	Type           string
	Location       string
	Name           string
	Output         OutputDevice
	Input          InputDevice
	Direction      string
	OnCommands     []string
	OffCommands    []string
	InitialValue   string
	targetValue    *Value
	UID            string
	lastCmd        string
	status         bool
	compare        Comparitor
	shutdown       bool
	filterMessages bool
	units          string
	Operator       string
	out            chan<- Message
	devIn          chan Message
	timerIn        chan bool
	timerOut       chan bool
}

//All gadgets respond to Robot Command Language (RCL) messages.  isMyCommand
//reads an RCL message and decides if it was meant for this instance
//of Gadget.
func (g *Gadget) isMyCommand(msg *Message) (bool, string, string) {
	if msg.Type != COMMAND {
		return false, "", ""
	}

	if msg.Body == "update" || msg.Body == "shutdown" {
		return true, "", ""
	}

	for _, cmd := range g.OnCommands {
		if strings.Index(msg.Body, cmd) == 0 {
			return true, "on", cmd
		}
	}

	for _, cmd := range g.OffCommands {
		if strings.Index(msg.Body, cmd) == 0 {
			return true, "off", cmd
		}
	}
	return false, "", ""
}

func (g *Gadget) GetDirection() string {
	if g.Output != nil {
		return "output"
	}
	return "input"
}

//Start is one of the two interface methods of GoGadget.  Start takes
//in in and out chan and is meant to be called as a goroutine.
func (g *Gadget) Start(in <-chan Message, out chan<- Message) {
	g.out = out
	g.timerIn = make(chan bool)
	g.timerOut = make(chan bool)
	if g.Output != nil {
		if len(g.InitialValue) > 0 {
			g.readInitialValue()
		} else {
			g.off()
		}
		g.doOutputLoop(in)
	} else if g.Input != nil {
		g.doInputLoop(in)
	}
}

//Once a gadget is started as a goroutine, this loop collects
//all the messages that are sent to this particular Gadget and
//responds accordingly.  This is the loop that is executed if
//this Gadget is an input Gadget
func (g *Gadget) doInputLoop(in <-chan Message) {
	devOut := make(chan Value, 10)
	g.devIn = make(chan Message, 10)
	go g.Input.Start(g.devIn, devOut)
	g.sendUpdate()
	for !g.shutdown {
		select {
		case msg := <-in:
			g.readMessage(&msg)
		case val := <-devOut:
			g.out <- Message{
				UUID:      GetUUID(),
				Sender:    g.UID,
				Type:      "update",
				Location:  g.Location,
				Name:      g.Name,
				Value:     val,
				Timestamp: time.Now().UTC(),
				Info: Info{
					Direction: g.Direction,
					Type:      g.Type,
				},
			}
		}
	}
}

func (g *Gadget) readInitialValue() {
	msg := &Message{
		UUID: GetUUID(),
		Body: g.InitialValue,
	}
	g.readCommand(msg, "on", g.InitialValue)
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
	err := g.Output.On(val)
	if err != nil {
		log.Println("on err", err)
	} else if !g.status {
		g.targetValue = val
		g.status = true
		g.sendUpdate()
	}
}

func (g *Gadget) off() {
	g.status = false
	g.targetValue = nil
	g.Output.Off()
	g.compare = nil
	g.sendUpdate()
}

func (g *Gadget) readMessage(msg *Message) {
	if g.devIn != nil {
		g.devIn <- *msg
	}
	mine, onoff, matched := g.isMyCommand(msg)
	if mine && msg.Type == COMMAND {
		g.lastCmd = matched
		g.readCommand(msg, onoff, matched)
	} else if g.status && msg.Type == UPDATE {
		g.readUpdate(msg)
	}
}

func (g *Gadget) readUpdate(msg *Message) {
	if g.status && g.compare != nil && g.compare(msg) {
		g.off()
	} else if g.status && (msg.Location == g.Location || !g.filterMessages) {
		if g.Output.Update(msg) {
			g.sendUpdate()
		}
	}
}

func (g *Gadget) readCommand(msg *Message, onoff, matched string) {
	if msg.Body == "shutdown" {
		g.shutdown = true
		g.off()
	} else if msg.Body == "update" {
		g.sendUpdate()
	} else if onoff == "on" {
		g.readOnCommand(msg, matched)
	} else if onoff == "off" {
		g.readOffCommand(msg)
	}
}

func (g *Gadget) readOnCommand(msg *Message, matched string) {
	var val *Value
	if len(strings.Trim(msg.Body, " ")) > len(matched) {
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
		return val, fmt.Errorf("could not parse %s", cmd)
	}
	gadget, ok := units[unit]

	if !ok {
		return nil, nil
	}

	val = &Value{
		Value: value,
		Units: unit,
		Cmd:   cmd,
	}

	if gadget == "time" {
		go g.startTimer(value, unit, g.timerIn, g.timerOut)
	} else if gadget == "volume" {
		g.setCompare(value, unit, gadget)
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

func (g *Gadget) getDuration(value float64, unit string) time.Duration {
	if unit == "minute" || unit == "minutes" {
		value *= 60.0
	} else if unit == "hour" || unit == "hours" {
		value *= 3600.0
	}
	return time.Duration(value * float64(time.Second))
}

func (g *Gadget) startTimer(value float64, unit string, in <-chan bool, out chan<- bool) {
	d := g.getDuration(value, unit)
	keepGoing := true
	for keepGoing {
		select {
		case <-in:
			keepGoing = false
		case <-time.After(d):
			keepGoing = false
			out <- true
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
	var value Value
	if g.Input != nil {
		value = *(g.Input.GetValue())
	} else {
		value = Value{
			Units:  g.units,
			Value:  g.status,
			Output: g.Output.Status(),
			Cmd:    g.lastCmd,
		}
	}
	g.out <- Message{
		UUID:        GetUUID(),
		Sender:      g.UID,
		Type:        UPDATE,
		Location:    g.Location,
		Name:        g.Name,
		Value:       value,
		TargetValue: g.targetValue,
		Timestamp:   time.Now().UTC(),
		Info: Info{
			Type:      g.Type,
			Direction: g.Direction,
			On:        g.OnCommands,
			Off:       g.OffCommands,
		},
	}
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
	if len(parts) != 2 {
		return "", "", fmt.Errorf("invalid command: %s", cmd)
	}
	return parts[0], parts[1], nil
}

func stripCommand(cmd string) string {
	cmd = strings.Trim(cmd, " ")
	i := strings.Index(cmd, " for ")
	if i != -1 {
		return cmd[i+5:]
	}
	i = strings.Index(cmd, " to ")
	if i != -1 {
		return cmd[i+4:]
	}
	return ""
}

// GetUUID generates a random UUID according to RFC 4122
func GetUUID() string {
	uuid := make([]byte, 16)
	n, err := io.ReadFull(rand.Reader, uuid)
	if n != len(uuid) || err != nil {
		return ""
	}
	// variant bits; see section 4.1.1
	uuid[8] = uuid[8]&^0xc0 | 0x80
	// version 4 (pseudo-random); see section 4.1.3
	uuid[6] = uuid[6]&^0xf0 | 0x40
	return fmt.Sprintf("%x-%x-%x-%x-%x", uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}
