package output

import (
	"fmt"
	"log"
	"strings"
	"strconv"
	"bitbucket.com/cswank/gogadgets"
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

type Comparitor func(x float64) bool

type Trigger struct {
	location string
	name string
	uid string
	operator string
	command string
	offCommand string
	triggerType string
	compare Comparitor
	in <-chan gogadgets.Message
	out chan<- gogadgets.Message
}

func (t *Trigger) Start(out chan<- gogadgets.Message, in <-chan gogadgets.Message) {
	t.uid = fmt.Sprintf("%s %s trigger", t.location, t.name)
	t.in = in
	t.out = out
	t.parseCommand()
}

func (t *Trigger) parseCommand() {
	value, unit, err := t.getValue()
	if err != nil {
		log.Println("could not parse", t.command)
	}
	triggerType, ok := units[unit]
	if ok {
		t.triggerType = triggerType
		if t.triggerType == "time" {
			t.waitForTime(value, unit)
		} else if t.triggerType == "volume" || t.triggerType == "temperature" {
			t.waitForMessage(value, unit)
		}
	}
}

func (t *Trigger) waitForMessage(value float64, unit string) {
	if t.operator == "<=" {
		t.compare = func(x float64) bool {return x <= value}
		t.doWaitForMessage()
	} else if t.operator == ">=" {
		t.compare = func(x float64) bool {return x >= value}
		t.doWaitForMessage()
	}
}

func (t *Trigger) getValue() (float64, string, error) {
	t.stripCommand()
	value, unit, err := t.splitCommand()
	var v float64
	if err == nil {
		v, err = strconv.ParseFloat(value, 64)
	}
	return v, unit, err
}

func (t *Trigger) splitCommand() (string, string, error) {
	parts := strings.Split(t.command, " ")
	return parts[0], parts[1], nil
}

func (t *Trigger) getMessage() gogadgets.Message {
	return gogadgets.Message{
		Sender: t.uid,
		Type: gogadgets.COMMAND,
		Body: t.offCommand,
	}
}

func (t *Trigger) doWaitForMessage() {
	msg := <-t.in
	if msg.Type == "command" && msg.Body == "stop" {
		return
	}
	val, ok := msg.Locations[t.location].Input[t.triggerType].Value.(float64)
	if ok && t.compare(val) {
		t.out<- t.getMessage()
	} else {
		t.doWaitForMessage()
	}
}

func (t *Trigger) waitForTime(value float64, unit string) {
	//msg := <-t.in
}

func (t *Trigger) stripCommand() {
	t.command = strings.Trim(t.command, " ")
	t.command = strings.TrimPrefix(t.command, "for ")
	t.command = strings.TrimPrefix(t.command, "to ")
}
