package output

import (
	"fmt"
	"time"
	"testing"
	"bitbucket.com/cswank/gogadgets/utils"
	"bitbucket.com/cswank/gogadgets/devices"
	"bitbucket.com/cswank/gogadgets/models"
)

type FakeOutput struct {
	devices.OutputDevice
	on bool
}

func (f *FakeOutput) Update(msg *models.Message) {
	
}

func (f *FakeOutput) On(val *models.Value) error {
	f.on = true
	return nil
}

func (f *FakeOutput) Off() error {
	f.on = false
	return nil
}

func (f *FakeOutput) Status() interface{} {
	return f.on
}

type FakePoller struct {
	devices.Poller
	val bool
}

func (f *FakePoller) Wait() (bool, error) {
	time.Sleep(100 * time.Millisecond)
	f.val = !f.val
	return f.val, nil
}


func TestStripCommand(t *testing.T) {
	tr := Gadget{
		location: "tank",
		name: "valve",
		operator: ">=",
		onCommand: "fill tank",
		offCommand: "stop filling tank",
	}
	cmd := tr.stripCommand("fill tank to 5 liters")
	if cmd != "5 liters" {
		t.Error(cmd)
	}
}

func TestGetValue(t *testing.T) {
	g := Gadget{
		location: "tank",
		name: "valve",
		operator: ">=",
		offCommand: "stop filling tank",
	}
	val, unit, err := g.getValue("5 liters")
	if err != nil {
		t.Error(err)
	}
	if val != 5.0 {
		t.Error("incorrect value", val)
	}
	if unit != "liters" {
		t.Error("incorrect unit", unit)
	}
}

func TestGetTimeValue(t *testing.T) {
	g := Gadget{
		location: "tank",
		name: "valve",
		operator: ">=",
		offCommand: "stop filling tank",
	}
	val, unit, err := g.getValue("1.1 minutes")
	if err != nil {
		t.Error(err)
	}
	if val != 1.1 {
		t.Error("incorrect value", val)
	}
	if unit != "minutes" {
		t.Error("incorrect unit", unit)
	}
}

func TestStart(t *testing.T) {
	location := "lab"
	name := "led"
	g := Gadget{
		location: location,
		name: name,
		onCommand: fmt.Sprintf("turn on %s %s", location, name),
		offCommand: fmt.Sprintf("turn off %s %s", location, name),
		output: &FakeOutput{},
		uid: fmt.Sprintf("%s %s", location, name),
	}
	input := make(chan models.Message)
	output := make(chan models.Message)
	go g.Start(input, output)
	msg := models.Message{
		Type: "command",
		Body: "turn on lab led",
	}
	input<- msg
	status := <-output
	if status.Value.Value != true {
		t.Fatal("shoulda been on", status)
	}
	
	msg = models.Message{
		Type: "command",
		Body: "turn off lab led",
	}
	input<- msg
	status = <-output
	if status.Value.Value != false {
		t.Error("shoulda been off", status)
	}
	msg = models.Message{
		Type: "command",
		Body: "shutdown",
	}
	input<- msg
	status = <-output
}

func TestStartWithTrigger(t *testing.T) {
	location := "tank"
	name := "valve"
	g := Gadget{
		location: location,
		name: name,
		operator: ">=",
		onCommand: fmt.Sprintf("fill %s", location),
		offCommand: fmt.Sprintf("stop filling %s", location),
		output: &FakeOutput{},
		uid: fmt.Sprintf("%s %s", location, name),
	}
	input := make(chan models.Message)
	output := make(chan models.Message)
	go g.Start(input, output)
	msg := models.Message{
		Type: "command",
		Body: "fill tank to 4.4 liters",
	}
	input<- msg
	status := <-output
	if status.Value.Value != true {
		t.Error("shoulda been on", status)
	}
	//make a message that should trigger the trigger and stop the device
	msg = models.Message{
		Sender: "tank volume",
		Type: models.STATUS,
		Location: "tank",
		Name: "volume",
		Value: models.Value{
			Units: "liters",
			Value: 4.4,
		},
	}
	input<- msg
	status = <-output
	if status.Value.Value != false {
		t.Error("shoulda been off", status)
	}
}

func TestStartWithTimeTrigger(t *testing.T) {
	location := "lab"
	name := "led"
	g := Gadget{
		location: location,
		name: name,
		onCommand: "turn on lab led",
		operator: ">=",
		offCommand: "turn off lab led",
		output: &FakeOutput{},
		uid: fmt.Sprintf("%s %s", location, name),
	}
	input := make(chan models.Message)
	output := make(chan models.Message)
	go g.Start(input, output)
	msg := models.Message{
		Type: "command",
		Body: "turn on lab led for 0.01 seconds",
	}
	input<- msg
	status := <-output
	if status.Value.Value != true {
		t.Error("shoulda been on", status)
	}
	//wait for a second
	status = <-output
	if status.Value.Value != false {
		t.Error("shoulda been off", status)
	}
}

func TestStartWithTimeTriggerWithInterrupt(t *testing.T) {
	location := "lab"
	name := "led"
	g := Gadget{
		location: location,
		name: name,
		onCommand: "turn on lab led",
		offCommand: "turn off lab led",
		output: &FakeOutput{},
		uid: fmt.Sprintf("%s %s", location, name),
	}
	input := make(chan models.Message)
	output := make(chan models.Message)
	go g.Start(input, output)
	msg := models.Message{
		Type: "command",
		Body: "turn on lab led for 30 seconds",
	}
	input<- msg
	status := <-output
	if status.Value.Value != true {
		t.Error("shoulda been on", status)
	}
	
	msg = models.Message{
		Type: "command",
		Body: "turn on lab led",
	}
	input<- msg

	msg = models.Message{
		Type: "status",
		Body: "",
	}
	input<- msg

	msg = models.Message{
		Type: "command",
		Body: "turn off lab led",
	}
	input<- msg
	status = <-output
	if status.Value.Value != false {
		t.Error("shoulda been off", status)
	}
}

func TestStartWithTimeTriggerForReals(t *testing.T) {
	if !utils.FileExists("/sys/class/gpio/export") {
		return //not a beaglebone
	}
	pin := &models.Pin{Type:"gpio", Port: "9", Pin: "15"}
	gpio, err := devices.NewGPIO(pin)
	if err != nil {
		t.Fatal(err)
	}
	location := "lab"
	name := "led"
	g := Gadget{
		location: location,
		name: name,
		onCommand: "turn on lab led",
		offCommand: "turn off lab led",
		output: gpio,
		uid: fmt.Sprintf("%s %s", location, name),
	}
	input := make(chan models.Message)
	output := make(chan models.Message)
	go g.Start(input, output)
	msg := models.Message{
		Type: "command",
		Body: "turn on lab led for 0.1 seconds",
	}
	input<- msg
	status := <-output
	if status.Value.Value != true {
		t.Error("shoulda been on", status)
	}
	//wait for a second
	status = <-output
	if status.Value.Value != false {
		t.Error("shoulda been off", status)
	}
}


func TestInputStart(t *testing.T) {
	location := "lab"
	name := "switch"
	poller := &FakePoller{}
	s := &devices.Switch{
		GPIO: poller,
	}
	g := Gadget{
		location: location,
		name: name,
		input: s,
		uid: fmt.Sprintf("%s %s", location, name),
	}
	input := make(chan models.Message)
	output := make(chan models.Message)
	go g.Start(input, output)
	val := <-output
	if val.Value.Value != true {
		t.Error("should have been true")
	}
	val = <-output
	if val.Value.Value != false {
		t.Error("should have been false")
	}
}
