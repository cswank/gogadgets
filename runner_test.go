package gogadgets

import (
	"testing"
	"time"
)

func TestReadWaitCommand(t *testing.T) {
	m := Runner{}
	waitTime, err := m.getWaitTime("wait for 3.3 seconds")
	if err != nil {
		t.Error(err)
	}
	if waitTime != time.Duration(3.3 * float64(time.Second)) {
		t.Error("incorrect time", waitTime)
	}
	waitTime, err = m.getWaitTime("wait for 1 second")
	if err != nil {
		t.Error(err)
	}
	if waitTime != time.Duration(1.0 * float64(time.Second)) {
		t.Error("incorrect time", waitTime)
	}
	waitTime, err = m.getWaitTime("wait for 10 hours")
	if err != nil {
		t.Error(err)
	}
	if waitTime != time.Duration(36000.0 * float64(time.Second)) {
		t.Error("incorrect time", waitTime)
	}
	waitTime, err = m.getWaitTime("wait for 1.1 minutes")
	if err != nil {
		t.Error(err)
	}
	if waitTime != time.Duration(66.0 * float64(time.Second)) {
		t.Error("incorrect time", waitTime)
	}
}

func TestStepExp(t *testing.T) {
	cmd := "wait for tank volume <= 5.4"
	result := stepExp.FindStringSubmatch(cmd)
	if len(result) != 4 {
		t.Fatal(result)
	}
	if result[3] != "5.4" {
		t.Error(result)
	}
	if result[2] != "<=" {
		t.Error(result)
	}
	if result[1] != "tank volume" {
		t.Error(result)
	}

	cmd = "wait for fish tank temperature > 31"
	result = stepExp.FindStringSubmatch(cmd)
	if len(result) != 4 {
		t.Fatal(result)
	}
	if result[2] != ">" {
		t.Error(result)
	}
	if result[3] != "31" {
		t.Error(result)
	}
	if result[1] != "fish tank temperature" {
		t.Error(result)
	}
}

func TestSetStepChecker(t *testing.T) {
	m := Runner{}
	cmd := "wait for tank volume >= 5.4"
	m.setStepChecker(cmd)
	msg := &Message{
		Sender: "tank volume",
		Value: Value{
			Value: 5.4,
		},
	}
	if !m.stepChecker(msg) {
		t.Error("should have been true")
	}

	msg = &Message{
		Sender: "fish tank volume",
		Value: Value{
			Value: 5.4,
		},
	}
	if m.stepChecker(msg) {
		t.Error("should have been false")
	}
}

func TestParseWaitCommand(t *testing.T) {
	m := Runner{}
	cmd := "wait for tank volume >= 5.4"
	uid, operator, value, err := m.parseWaitCommand(cmd)
	if err != nil {
		t.Error(err)
	}
	if value != 5.4 {
		t.Error(value)
	}
	if uid != "tank volume" {
		t.Error(uid)
	}
	if operator != ">=" {
		t.Error(operator)
	}
	cmd = "wait for fish tank temperature > 31"
	uid, operator, value, err = m.parseWaitCommand(cmd)
	if err != nil {
		t.Fatal(err)
	}
	if value != 31.0 {
		t.Error(value)
	}
	if uid != "fish tank temperature" {
		t.Error(uid)
	}
	if operator != ">" {
		t.Error(operator)
	}
}

func TestRunMethod(t *testing.T) {
	in := make(chan Message)
	out := make(chan Message)
	m := Runner{}
	go m.Start(out, in)
	msg := Message{
		Type: METHOD,
		Method: Method{
			Steps: []string{
				"fill boiler to 3.3 gallons",
				"heat boiler to 95 C",
				"wait for boiler temperature >= 95",
				"stop heating boiler",
			},
		},
	}
	out<- msg
	<-in
	msg = <-in
	if msg.Type != "command" && msg.Body != "fill boiler to 3.3 gallons" {
		t.Error(msg)
	}
	<-in
	msg = <-in
	if msg.Type != "command" && msg.Body != "heat boiler to 95 C" {
		t.Error(msg)
	}
	msg = Message{
		Type: "update",
		Sender: "boiler temperature",
		Value: Value{
			Value: 96.0,
			Units: "C",
		},
	}
	<-in
	out<- msg
	<-in
	msg = <-in
	if msg.Type != "command" && msg.Body != "stop heating boiler" {
		t.Error(msg)
	}
	msg = Message{
		Type: "command",
		Body: "shutdown",
	}
	<-in
	out<- msg
	<-in
}

func TestUserStepChecker(t *testing.T) {
	m := Runner{}
	m.setUserStepChecker("wait for user to laugh")
	msg := &Message{
		Type: "update",
		Body: "wait for user to cry",
	}
	if m.stepChecker(msg) {
		t.Error("should have returned false")
	}
	msg.Body = "wait for user to laugh"
	if !m.stepChecker(msg) {
		t.Error("should have returned true")
	}
}

func TestRunAnotherMethod(t *testing.T) {
	in := make(chan Message)
	out := make(chan Message)
	m := Runner{}
	go m.Start(out, in)
	msg := Message{
		Type: METHOD,
		Method: Method{
			Steps:[]string{
				"turn on lab led",
				"wait for 0.1 seconds",
				"turn off lab led",
				"wait for user to turn off power",
				"shutdown",
			},
		},
	}
	out<- msg
	<-in
	msg = <-in
	if msg.Type != "command" && msg.Body != "turn on lab led" {
		t.Error(msg)
	}
	<-in
	<-in
	<-in
	<-in
	msg = <-in
	if msg.Type != "command" && msg.Body != "turn off lab led" {
		t.Error(msg)
	}
	<-in
	msg = <-in
	if msg.Type != "command" && msg.Body != "wait for user to turn off power" {
		t.Error(msg)
	}
	out<- Message{
		Type: "update",
		Body: "wait for user to turn off power",
	}
	<-in
	msg = <-in
	if msg.Type != "command" && msg.Body != "shutdown" {
		t.Error(msg)
	}
}



