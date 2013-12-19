package methods

import (
	"fmt"
	"errors"
	"time"
	"regexp"
	"strings"
	"strconv"
	"bitbucket.com/cswank/gogadgets/models"
)

var (
	timeExp = regexp.MustCompile(`for (\d*\.?\d*) (second|seconds|minute|minutes|hour|hours)`)
)

type stepChecker func(msg *models.Message) bool

type Methods struct {
	models.Gadget
	method []string
	waitTime time.Duration
	stepChecker stepChecker
	step int
	out chan<- models.Message
}

func (m *Methods) Start(in <-chan models.Message, out chan<- models.Message) {
	shutdown := false
	m.out = out
	for !shutdown {
		select {
		case msg := <-in:
			m.readMessage(&msg)
		case <-time.After(m.waitTime):
			
		}
	}
}

func (m *Methods) readMessage(msg *models.Message) {
	if msg.Type == models.METHOD {
		m.method = msg.Method
		m.step = -1
		m.runNextStep()
	} else if len(m.method) != 0 && msg.Type == models.UPDATE {
		m.checkUpdate(msg)
	}
}

func (m *Methods) checkUpdate(msg *models.Message) {
	if m.stepChecker != nil && m.stepChecker(msg) {
		m.runNextStep()
	}
}

func (m *Methods) runNextStep() {
	m.step += 1
	if len(m.method) <= m.step {
		m.method = []string{}
		m.step = -1
		return
	}
	cmd := m.method[m.step]
	if strings.Index(cmd, "wait") == 0 {
		m.readWaitCommand(cmd)
	} else {
		m.sendCommand(cmd)
		m.runNextStep()
	}
}

func (m *Methods) sendCommand(cmd string) {
	msg := models.Message{
		Type: models.COMMAND,
		Body: cmd,
	}
	m.out<- msg
}

func (m *Methods) readWaitCommand(cmd string) (d time.Duration, err error) {
	result := timeExp.FindStringSubmatch(cmd)
	if len(result) == 3 {
		units := result[2]
		t, err := strconv.ParseFloat(result[1], 64)
		if err != nil {
			err = errors.New(fmt.Sprintf("could not parse command", cmd))
		} else {
			if units == "minutes" || units == "minute" {
				t *= 60.0
			} else if units == "hours" || units == "hour" {
				t *= 3600.0
			}
			d = time.Duration(t * float64(time.Second))
		}
	} else {
		err = errors.New(fmt.Sprintf("could not parse command", cmd))
	}
	return d, err
}
