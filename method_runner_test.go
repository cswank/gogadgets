package gogadgets

// import (
// 	"testing"
// 	"time"
// )

// func TestReadWaitCommand(t *testing.T) {
// 	m := MethodRunner{}
// 	waitTime, err := m.getWaitTime("wait for 3.3 seconds")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if waitTime != time.Duration(3.3*float64(time.Second)) {
// 		t.Error("incorrect time", waitTime)
// 	}
// 	waitTime, err = m.getWaitTime("wait for 1 second")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if waitTime != time.Duration(1.0*float64(time.Second)) {
// 		t.Error("incorrect time", waitTime)
// 	}
// 	waitTime, err = m.getWaitTime("wait for 10 hours")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if waitTime != time.Duration(36000.0*float64(time.Second)) {
// 		t.Error("incorrect time", waitTime)
// 	}
// 	waitTime, err = m.getWaitTime("wait for 1.1 minutes")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if waitTime != time.Duration(66.0*float64(time.Second)) {
// 		t.Error("incorrect time", waitTime)
// 	}
// }

// func TestStepExp(t *testing.T) {
// 	cmd := "wait for tank volume <= 5.4 liters"
// 	result := stepExp.FindStringSubmatch(cmd)
// 	if len(result) != 5 {
// 		t.Fatal(result)
// 	}
// 	if result[4] != "liters" {
// 		t.Error(result)
// 	}
// 	if result[3] != "5.4" {
// 		t.Error(result)
// 	}
// 	if result[2] != "<=" {
// 		t.Error(result)
// 	}
// 	if result[1] != "tank volume" {
// 		t.Error(result)
// 	}

// 	cmd = "wait for fish tank temperature > 31 C"
// 	result = stepExp.FindStringSubmatch(cmd)
// 	if len(result) != 5 {
// 		t.Fatal(result)
// 	}
// 	if result[2] != ">" {
// 		t.Error(result)
// 	}
// 	if result[3] != "31" {
// 		t.Error(result)
// 	}
// 	if result[1] != "fish tank temperature" {
// 		t.Error(result)
// 	}
// }

// func TestSetStepChecker(t *testing.T) {
// 	m := MethodRunner{}
// 	cmd := "wait for tank volume >= 5.4 gallons"
// 	m.setStepChecker(cmd)
// 	msg := &Message{
// 		Sender: "tank volume",
// 		Value: Value{
// 			Value: 5.4,
// 		},
// 	}
// 	if !m.stepChecker(msg) {
// 		t.Error("should have been true")
// 	}

// 	msg = &Message{
// 		Sender: "fish tank volume",
// 		Value: Value{
// 			Value: 5.4,
// 		},
// 	}
// 	if m.stepChecker(msg) {
// 		t.Error("should have been false")
// 	}
// }

// func TestSetBoolStepChecker(t *testing.T) {
// 	m := MethodRunner{}
// 	cmd := "wait for lab switch == true"
// 	m.setStepChecker(cmd)
// 	msg := &Message{
// 		Sender: "lab switch",
// 		Value: Value{
// 			Value: true,
// 		},
// 	}
// 	if !m.stepChecker(msg) {
// 		t.Error("should have been true")
// 	}

// 	msg = &Message{
// 		Sender: "fish tank volume",
// 		Value: Value{
// 			Value: 5.4,
// 		},
// 	}
// 	if m.stepChecker(msg) {
// 		t.Error("should have been false")
// 	}
// }

// func TestParseWaitCommand(t *testing.T) {
// 	m := MethodRunner{}
// 	cmd := "wait for tank volume >= 5.4 gallons"
// 	uid, operator, value, err := m.parseWaitCommand(cmd)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if value != 5.4 {
// 		t.Error(value)
// 	}
// 	if uid != "tank volume" {
// 		t.Error(uid)
// 	}
// 	if operator != ">=" {
// 		t.Error(operator)
// 	}
// 	cmd = "wait for fish tank temperature > 31 C"
// 	uid, operator, value, err = m.parseWaitCommand(cmd)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if value != 31.0 {
// 		t.Error(value)
// 	}
// 	if uid != "fish tank temperature" {
// 		t.Error(uid)
// 	}
// 	if operator != ">" {
// 		t.Error(operator)
// 	}
// }

// func TestRunMethod(t *testing.T) {
// 	in := make(chan Message)
// 	out := make(chan Message)
// 	m := MethodRunner{}
// 	go m.Start(out, in)
// 	msg := Message{
// 		Type: METHOD,
// 		Method: Method{
// 			Steps: []string{
// 				"fill boiler to 3.3 gallons",
// 				"heat boiler to 95 C",
// 				"wait for boiler temperature >= 95 C",
// 				"stop heating boiler",
// 			},
// 		},
// 	}
// 	out <- msg
// 	<-in
// 	msg = <-in
// 	if msg.Type != "command" && msg.Body != "fill boiler to 3.3 gallons" {
// 		t.Error(msg)
// 	}
// 	<-in
// 	msg = <-in
// 	if msg.Type != "command" && msg.Body != "heat boiler to 95 C" {
// 		t.Error(msg)
// 	}
// 	msg = Message{
// 		Type:   "update",
// 		Sender: "boiler temperature",
// 		Value: Value{
// 			Value: 96.0,
// 			Units: "C",
// 		},
// 	}
// 	<-in
// 	out <- msg
// 	<-in
// 	msg = <-in
// 	if msg.Type != "command" && msg.Body != "stop heating boiler" {
// 		t.Error(msg)
// 	}
// 	msg = Message{
// 		Type: "command",
// 		Body: "shutdown",
// 	}
// 	<-in
// 	out <- msg
// 	<-in
// }

// func TestUserStepChecker(t *testing.T) {
// 	m := MethodRunner{}
// 	m.setUserStepChecker("wait for user to laugh")
// 	msg := &Message{
// 		Type: "update",
// 		Body: "wait for user to cry",
// 	}
// 	if m.stepChecker(msg) {
// 		t.Error("should have returned false")
// 	}
// 	msg.Body = "wait for user to laugh"
// 	if !m.stepChecker(msg) {
// 		t.Error("should have returned true")
// 	}
// }

// func TestRunAnotherMethod(t *testing.T) {
// 	in := make(chan Message)
// 	out := make(chan Message)
// 	m := MethodRunner{}
// 	go m.Start(out, in)
// 	msg := Message{
// 		Type: METHOD,
// 		Method: Method{
// 			Steps: []string{
// 				"turn on lab led",
// 				"wait for 0.1 seconds",
// 				"turn off lab led",
// 				"wait for user to turn off power",
// 				"shutdown",
// 			},
// 		},
// 	}
// 	out <- msg
// 	msg = <-in
// 	if msg.Type != "method update" || msg.Method.Step != 0 {
// 		t.Error(msg)
// 	}
// 	msg = <-in
// 	if msg.Type != "command" && msg.Body != "turn on lab led" {
// 		t.Error(msg)
// 	}
// 	msg = <-in
// 	if msg.Type != "method update" || msg.Method.Step != 1 {
// 		t.Error(msg)
// 	}
// 	msg = <-in
// 	if msg.Type != "method update" || msg.Method.Step != 1 {
// 		t.Error(msg)
// 	}
// 	msg = <-in
// 	if msg.Type != "method update" || msg.Method.Step != 1 {
// 		t.Error(msg)
// 	}
// 	<-in
// 	msg = <-in
// 	if msg.Type != "command" || msg.Body != "turn off lab led" {
// 		t.Error(msg)
// 	}
// 	msg = <-in
// 	if msg.Type != "method update" || msg.Method.Step != 3 {
// 		t.Error(msg)
// 	}
// 	out <- Message{
// 		Type: "update",
// 		Body: "wait for user to turn off power",
// 	}
// 	msg = <-in
// 	if msg.Type != "method update" || msg.Method.Step != 4 {
// 		t.Error(msg)
// 	}
// 	msg = <-in
// 	if msg.Type != "command" || msg.Body != "shutdown" {
// 		t.Error(msg)
// 	}
// }

// func TestRunBrewMethod(t *testing.T) {
// 	in := make(chan Message)
// 	out := make(chan Message)
// 	m := MethodRunner{}
// 	go m.Start(out, in)
// 	msg := Message{
// 		Type: METHOD,
// 		Method: Method{
// 			Steps: []string{
// 				"fill hlt to 7.0 gallons",
// 				"heat hlt to 175.000000 F",
// 				"wait for hlt temperature >= 175.000000",
// 				"fill tun to 3.125000 gallons",
// 				"wait for tun volume >= 3.125000 gallons",
// 				"wait for user to add grains",
// 				"fill hlt to 1.0 gallons",
// 				"heat hlt to 185 F",
// 				"wait for 1 second",
// 				"fill tun to 4.125000 gallons",
// 				"wait for 1 second",
// 				"wait for user ready to recirculate",
// 				"fill boiler",
// 				"wait for user recirculated",
// 				"fill boiler to 4.125000 gallons",
// 				"heat boiler to 190 F",
// 				"wait for boiler volume >= 4.125000 gallons",
// 				"fill tun 3.125000 gallons",
// 				"wait for tun volume >= 3.125000 gallons",
// 				"stop heating hlt",
// 				"wait for 2 minutes",
// 				"wait for user ready to recirculate",
// 				"fill boiler",
// 				"wait for user recirculated",
// 				"fill boiler to 7.250000 gallons",
// 				"heat boiler to 204 F",
// 				"turn on fan",
// 				"wait for 60.000000 minutes",
// 				"stop heating boiler",
// 				"turn off fan",
// 				"wait for 5 minutes",
// 				"cool boiler to 80 F",
// 				"wait for boiler temperature <= 80 F",
// 				"wait for user to open ball valve",
// 				"fill fermenter",
// 				"wait for user to confirm boiler empty",
// 				"stop filling fermenter",
// 			},
// 		},
// 	}
// 	out <- msg
// 	msg = <-in
// 	if msg.Type != "method update" || msg.Method.Step != 0 {
// 		t.Error(msg)
// 	}
// 	msg = <-in
// 	if msg.Type != "command" && msg.Body != "fill hlt to 7.0 gallons" {
// 		t.Error(msg)
// 	}
// 	msg = <-in
// 	if msg.Type != "method update" || msg.Method.Step != 1 {
// 		t.Error(msg)
// 	}
// 	msg = <-in
// 	if msg.Type != "command" && msg.Body != "heat hlt to 175.000000 F" {
// 		t.Error(msg)
// 	}
// 	msg = <-in
// 	if msg.Type != "method update" || msg.Method.Step != 2 {
// 		t.Error(msg)
// 	}
// 	out <- Message{
// 		Type:   "update",
// 		Sender: "hlt temperature",
// 		Value:  Value{Value: 175.0, Units: "F"},
// 	}
// 	msg = <-in
// 	if msg.Type != "method update" || msg.Method.Step != 3 {
// 		t.Error(msg)
// 	}
// 	msg = <-in
// 	if msg.Type != "command" || msg.Body != "fill tun to 3.125000 gallons" {
// 		t.Error(msg)
// 	}
// 	msg = <-in
// 	if msg.Type != "method update" || msg.Method.Step != 4 {
// 		t.Error(msg)
// 	}
// 	out <- Message{
// 		Type:   "update",
// 		Sender: "tun volume",
// 		Value:  Value{Value: 3.125, Units: "gallons"},
// 	}
// 	msg = <-in
// 	if msg.Type != "method update" || msg.Method.Step != 5 {
// 		t.Error(msg)
// 	}
// 	out <- Message{
// 		Type: "update",
// 		Body: "wait for user to add grains",
// 	}
// 	msg = <-in
// 	if msg.Type != "method update" || msg.Method.Step != 6 {
// 		t.Error(msg)
// 	}
// 	msg = <-in
// 	if msg.Type != "command" || msg.Body != "fill hlt to 1.0 gallons" {
// 		t.Error(msg)
// 	}
// 	msg = <-in
// 	if msg.Type != "method update" || msg.Method.Step != 7 {
// 		t.Error(msg)
// 	}
// 	msg = <-in
// 	if msg.Type != "command" || msg.Body != "heat hlt to 185 F" {
// 		t.Error(msg)
// 	}
// 	for {
// 		msg = <-in
// 		if msg.Type != "method update" {
// 			break
// 		}
// 	}
// 	if msg.Type != "command" || msg.Body != "fill tun to 4.125000 gallons" {
// 		t.Error(msg)
// 	}
// 	for {
// 		msg = <-in
// 		if msg.Method.Step != 10 {
// 			break
// 		}
// 	}
// 	if msg.Type != "method update" || msg.Method.Step != 11 {
// 		t.Error(msg)
// 	}
// }
