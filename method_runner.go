package gogadgets

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	timeExp = regexp.MustCompile(`for (\d*\.?\d*) (seconds?|minutes?|hours?)`)
	stepExp = regexp.MustCompile(`for (.+) (>=|>|==|<=|<) ([\w\.]+) ?(.+)?`)
)

type stepChecker func(msg *Message) bool
type comparitor func(value interface{}) bool

//  Gadgets respond to the Robot Command Language (RCL) and a
//  list of RCL messages can be run to form a method.  Runner
//  takes a method as input and runs it.
//
//  If a RCL message starts with 'wait' runner pauses the method and
//  waits for the condition if the wait to be fulfilled.  For example,
//  if the RCL message
//      'wait for 5 seconds'
//  is recieved, Runner waits for 5 seconds and continues with the the
//  rest of the message.
//
//  Another example would be
//      'wait for boiler temperature >= 200 F'.
//  MethodRunner would then wait for a message from the boiler that says
//  its temperature is 200 F (or more).  It then sends the next
//  message of the method
type MethodRunner struct {
	method      Method
	stepChecker stepChecker
	step        int
	uid         string
	out         chan<- Message
	timeOut     chan bool
}

func (m *MethodRunner) GetUID() string {
	return "method runner"
}

func (m *MethodRunner) Start(in <-chan Message, out chan<- Message) {
	m.uid = m.GetUID()
	m.out = out
	shutdown := false
	m.timeOut = make(chan bool)
	for !shutdown {
		select {
		case msg := <-in:
			shutdown = m.readMessage(&msg)
		case <-m.timeOut:
			m.runNextStep()
		}
	}
	m.out <- Message{}
}

func (m *MethodRunner) readMessage(msg *Message) (shutdown bool) {
	if msg.Type == METHOD {
		m.method = msg.Method
		m.step = -1
		m.runNextStep()
		shutdown = false
	} else if msg.Type == COMMAND && msg.Body == "update" {
		m.sendUpdate()
	} else if msg.Type == COMMAND && msg.Body == "clear method" {
		m.clear()
		m.sendUpdate()
	} else if len(m.method.Steps) != 0 && msg.Type == UPDATE {
		m.checkUpdate(msg)
		shutdown = false
	} else if msg.Type == COMMAND && msg.Body == "shutdown" {
		shutdown = true
	} else {
		shutdown = false
	}
	return shutdown
}

func (m *MethodRunner) sendUpdate() {
	m.method.Step = m.step
	msg := Message{
		Sender: m.GetUID(),
		Type:   UPDATE,
		Method: m.method,
	}
	m.out <- msg
}

func (m *MethodRunner) checkUpdate(msg *Message) {
	if m.stepChecker != nil && m.stepChecker(msg) {
		m.stepChecker = nil
		m.runNextStep()
	}
}

func (m *MethodRunner) clear() {
	m.method = Method{}
	m.step = -1
}

func (m *MethodRunner) runNextStep() {
	m.step += 1
	m.out <- Message{
		Sender: m.uid,
		Type:   METHODUPDATE,
		Method: Method{
			Step: m.step,
		},
	}
	if len(m.method.Steps) <= m.step {
		m.clear()
		return
	}
	cmd := m.method.Steps[m.step]
	if strings.Index(cmd, "wait") == 0 {
		m.readWaitCommand(cmd)
	} else {
		m.sendCommand(cmd)
		m.runNextStep()
	}
}

func (m *MethodRunner) sendCommand(cmd string) {
	msg := Message{
		Sender: m.uid,
		Type:   COMMAND,
		Body:   cmd,
	}
	m.out <- msg
}

func (m *MethodRunner) readWaitCommand(cmd string) {
	waitTime, err := m.getWaitTime(cmd)
	if strings.Index(cmd, "wait for user") == 0 {
		m.setUserStepChecker(cmd)
	} else if err == nil {
		go m.doCountdown(waitTime)
	} else {
		m.setStepChecker(cmd)
	}
}

func (m *MethodRunner) setUserStepChecker(cmd string) {
	m.stepChecker = func(msg *Message) bool {
		return msg.Body == cmd
	}
}

func (m *MethodRunner) setStepChecker(cmd string) {
	uid, operator, value, err := m.parseWaitCommand(cmd)
	if err == nil {
		compare, err := m.getCompare(operator, value)
		if err == nil {
			m.stepChecker = func(msg *Message) bool {
				return msg.Sender == uid &&
					compare(msg.Value.Value)
			}
		} else {
			log.Println(err)
			m.stepChecker = func(msg *Message) bool {
				return false
			}
		}
	} else {
		log.Println(err)
	}
}

func (m *MethodRunner) getCompare(operator string, value interface{}) (cmp comparitor, err error) {
	switch v := value.(type) {
	case float64:
		return m.getFloatCompare(operator, v)
	default:
		return func(x interface{}) bool { return x == v }, nil
	}
}

func (m *MethodRunner) getFloatCompare(operator string, value float64) (cmp comparitor, err error) {
	if operator == "<=" {
		cmp = func(x interface{}) bool { return x.(float64) <= value }
	} else if operator == "<" {
		cmp = func(x interface{}) bool { return x.(float64) < value }
	} else if operator == "==" {
		cmp = func(x interface{}) bool { return x.(float64) == value }
	} else if operator == ">=" {
		cmp = func(x interface{}) bool { return x.(float64) >= value }
	} else if operator == ">" {
		cmp = func(x interface{}) bool { return x.(float64) > value }
	} else {
		err = errors.New(fmt.Sprintf("invalid operator: %s", operator))
	}
	return cmp, err
}

func (m *MethodRunner) parseWaitCommand(cmd string) (uid, operator string, value interface{}, err error) {
	result := stepExp.FindStringSubmatch(cmd)
	if len(result) == 5 {
		uid = result[1]
		operator = result[2]
		v := result[3]
		if v == "true" {
			value = true
		} else if v == "false" {
			value = false
		} else {
			value, err = strconv.ParseFloat(v, 64)
		}
	}
	return uid, operator, value, err
}

func (m *MethodRunner) getWaitTime(cmd string) (waitTime time.Duration, err error) {
	result := timeExp.FindStringSubmatch(cmd)
	if len(result) != 3 {
		err = errors.New(fmt.Sprintf("could not parse command %s", cmd))
		return waitTime, err
	}
	units := result[2]
	t, err := strconv.ParseFloat(result[1], 64)
	if err != nil {
		err = errors.New(fmt.Sprintf("could not parse command %s", cmd))
		return waitTime, err
	} else {
		if units == "minutes" || units == "minute" {
			t *= 60.0
		} else if units == "hours" || units == "hour" {
			t *= 3600.0
		}
		waitTime = time.Duration(t * float64(time.Second))
	}
	return waitTime, err
}

func (m *MethodRunner) doCountdown(waitTime time.Duration) {
	t1 := time.Now()
	sleepTime := time.Duration(1 * time.Second)
	i := 0.0
	m.out <- Message{
		Sender: m.uid,
		Type:   METHODUPDATE,
		Method: Method{
			Time: int(waitTime.Seconds()),
			Step: m.step,
		},
	}
	for {
		time.Sleep(sleepTime)
		i += 1.0
		t2 := time.Now()
		d := t2.Sub(t1)
		sleepTime = time.Duration((1 - (d.Seconds() - i)) * float64(time.Second))
		m.out <- Message{
			Sender: m.uid,
			Type:   METHODUPDATE,
			Method: Method{
				Time: int(1 + waitTime.Seconds() - d.Seconds()),
				Step: m.step,
			},
		}
		if d > waitTime {
			m.timeOut <- true
			return
		}
	}
}
