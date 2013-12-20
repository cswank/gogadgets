package methods

import (
	"testing"
	"time"
	"bitbucket.com/cswank/gogadgets/models"
)

func TestReadWaitCommand(t *testing.T) {
	m := Methods{}
	m.readWaitCommand("wait for 3.3 seconds")
	if m.waitTime != time.Duration(3.3 * float64(time.Second)) {
		t.Error("incorrect time", m.waitTime)
	}
	m.readWaitCommand("wait for 1 second")
	if m.waitTime != time.Duration(1.0 * float64(time.Second)) {
		t.Error("incorrect time", m.waitTime)
	}
	m.readWaitCommand("wait for 10 hours")
	if m.waitTime != time.Duration(36000.0 * float64(time.Second)) {
		t.Error("incorrect time", m.waitTime)
	}
	m.readWaitCommand("wait for 1.1 minutes")
	if m.waitTime != time.Duration(66.0 * float64(time.Second)) {
		t.Error("incorrect time", m.waitTime)
	}
}

func TestStepExp(t *testing.T) {
	cmd := "wait for tank volume is 5.4"
	result := stepExp.FindStringSubmatch(cmd)
	if len(result) != 3 {
		t.Fatal(result)
	}
	if result[2] != "5.4" {
		t.Error(result)
	}
	if result[1] != "tank volume" {
		t.Error(result)
	}

	cmd = "wait for fish tank temperature is 31"
	result = stepExp.FindStringSubmatch(cmd)
	if len(result) != 3 {
		t.Fatal(result)
	}
	if result[2] != "31" {
		t.Error(result)
	}
	if result[1] != "fish tank temperature" {
		t.Error(result)
	}
}

func TestSetStepChecker(t *testing.T) {
	m := Methods{}
	cmd := "wait for tank volume is 5.4"
	m.setStepChecker(cmd)
	msg := &models.Message{
		Sender: "tank volume",
		Value: models.Value{
			Value: 5.4,
		},
	}
	if !m.stepChecker(msg) {
		t.Error("should have been true")
	}

	msg = &models.Message{
		Sender: "fish tank volume",
		Value: models.Value{
			Value: 5.4,
		},
	}
	if m.stepChecker(msg) {
		t.Error("should have been false")
	}
}

func TestParseWaitCommand(t *testing.T) {
	m := Methods{}
	cmd := "wait for tank volume is 5.4"
	uid, value, err := m.parseWaitCommand(cmd)
	if err != nil {
		t.Error(err)
	}
	if value != 5.4 {
		t.Error(value)
	}
	if uid != "tank volume" {
		t.Error(uid)
	}
	cmd = "wait for fish tank temperature is 31"
	uid, value, err = m.parseWaitCommand(cmd)
	if err != nil {
		t.Fatal(err)
	}
	if value != 31.0 {
		t.Error(value)
	}
	if uid != "fish tank temperature" {
		t.Error(uid)
	}
}

func TestRunMethod(t *testing.T) {
	in := make(chan models.Message)
	out := make(chan models.Message)
	m := Methods{}
	go m.Start(out, in)
	msg := models.Message{
		Type: models.METHOD,
		Method: []string{
			"fill boiler to 3.3 gallons",
			"heat boiler to 95 C",
			"wait for boiler temperature is 95",
			"stop heating boiler",
		},
	}
	out<- msg
	msg = <-in
	if msg.Type != "command" && msg.Body != "fill boiler to 3.3 gallons" {
		t.Error(msg)
	}
	msg = <-in
	if msg.Type != "command" && msg.Body != "heat boiler to 95 C" {
		t.Error(msg)
	}
	msg = models.Message{
		Type: "update",
		Sender: "boiler temperature",
		Value: models.Value{
			Value: 96.0,
			Units: "C",
		},
	}
	out<- msg
	msg = <-in
	if msg.Type != "command" && msg.Body != "stop heating boiler" {
		t.Error(msg)
	}
	msg = models.Message{
		Type: "command",
		Body: "shutdown",
	}
	out<- msg
	<-in
}

func TestRunAnotherMethod(t *testing.T) {
	in := make(chan models.Message)
	out := make(chan models.Message)
	m := Methods{}
	go m.Start(out, in)
	msg := models.Message{
		Type: models.METHOD,
		Method: []string{
			"turn on lab led",
			"wait for 0.1 seconds",
			"turn off lab led",
			"shutdown",
		},
	}
	out<- msg
	msg = <-in
	if msg.Type != "command" && msg.Body != "turn on lab led" {
		t.Error(msg)
	}
	msg = <-in
	if msg.Type != "command" && msg.Body != "turn off lab led" {
		t.Error(msg)
	}
	msg = <-in
	if msg.Type != "command" && msg.Body != "shutdown" {
		t.Error(msg)
	}
}



