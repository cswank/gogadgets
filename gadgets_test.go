package gogadgets

import (
	"fmt"
	"time"
	"testing"
	"bitbucket.com/cswank/gogadgets/utils"
)

func TestStripCommand(t *testing.T) {
	tr := Gadget{
		Location: "tank",
		Name: "valve",
		Operator: ">=",
		OnCommand: "fill tank",
		OffCommand: "stop filling tank",
	}
	cmd := tr.stripCommand("fill tank to 5 liters")
	if cmd != "5 liters" {
		t.Error(cmd)
	}
}

func TestGetValue(t *testing.T) {
	g := Gadget{
		Location: "tank",
		Name: "valve",
		Operator: ">=",
		OffCommand: "stop filling tank",
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
		Location: "tank",
		Name: "valve",
		Operator: ">=",
		OffCommand: "stop filling tank",
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
		Location: location,
		Name: name,
		OnCommand: fmt.Sprintf("turn on %s %s", location, name),
		OffCommand: fmt.Sprintf("turn off %s %s", location, name),
		Output: &FakeOutput{},
		UID: fmt.Sprintf("%s %s", location, name),
	}
	input := make(chan Message)
	output := make(chan Message)
	go g.Start(input, output)
	msg := Message{
		Type: "command",
		Body: "turn on lab led",
	}
	input<- msg
	status := <-output
	if status.Value.Value != true {
		t.Fatal("shoulda been on", status)
	}
	
	msg = Message{
		Type: "command",
		Body: "turn off lab led",
	}
	input<- msg
	status = <-output
	if status.Value.Value != false {
		t.Error("shoulda been off", status)
	}
	msg = Message{
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
		Location: location,
		Name: name,
		Operator: ">=",
		OnCommand: fmt.Sprintf("fill %s", location),
		OffCommand: fmt.Sprintf("stop filling %s", location),
		Output: &FakeOutput{},
		UID: fmt.Sprintf("%s %s", location, name),
	}
	input := make(chan Message)
	output := make(chan Message)
	go g.Start(input, output)
	msg := Message{
		Type: "command",
		Body: "fill tank to 4.4 liters",
	}
	input<- msg
	status := <-output
	if status.Value.Value != true {
		t.Error("shoulda been on", status)
	}
	//make a message that should trigger the trigger and stop the device
	msg = Message{
		Sender: "tank volume",
		Type: STATUS,
		Location: "tank",
		Name: "volume",
		Value: Value{
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
		Location: location,
		Name: name,
		OnCommand: "turn on lab led",
		Operator: ">=",
		OffCommand: "turn off lab led",
		Output: &FakeOutput{},
		UID: fmt.Sprintf("%s %s", location, name),
	}
	input := make(chan Message)
	output := make(chan Message)
	go g.Start(input, output)
	msg := Message{
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
		Location: location,
		Name: name,
		OnCommand: "turn on lab led",
		OffCommand: "turn off lab led",
		Output: &FakeOutput{},
		UID: fmt.Sprintf("%s %s", location, name),
	}
	input := make(chan Message)
	output := make(chan Message)
	go g.Start(input, output)
	msg := Message{
		Type: "command",
		Body: "turn on lab led for 30 seconds",
	}
	input<- msg
	status := <-output
	if status.Value.Value != true {
		t.Error("shoulda been on", status)
	}
	
	msg = Message{
		Type: "command",
		Body: "turn on lab led",
	}
	input<- msg

	msg = Message{
		Type: "status",
		Body: "",
	}
	input<- msg

	msg = Message{
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
	pin := &Pin{Type:"gpio", Port: "9", Pin: "15"}
	gpio, err := NewGPIO(pin)
	if err != nil {
		t.Fatal(err)
	}
	location := "lab"
	name := "led"
	g := Gadget{
		Location: location,
		Name: name,
		OnCommand: "turn on lab led",
		OffCommand: "turn off lab led",
		Output: gpio,
		UID: fmt.Sprintf("%s %s", location, name),
	}
	input := make(chan Message)
	output := make(chan Message)
	go g.Start(input, output)
	msg := Message{
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

func TestRealInput(t *testing.T) {
	if !utils.FileExists("/sys/class/gpio/export") {
		return //not a beaglebone
	}

	gpioConfig := &Pin{
		Type: "gpio",
		Port: "9",
		Pin: "15",
	}
	gpio, err := NewGPIO(gpioConfig)
	if err != nil {
		t.Fatal(err)
	}
	
	config := &Config{
		Location: "lab",
		Name: "switch",
		Pin: Pin{
			Type: "switch",
			Port: "9",
			Pin: "16",
			Edge: "both",
			Direction: "in",
			Value: 5.0,
			Units: "liters",
		},
	}
	s, err := NewGadget(config)
	if err != nil {
		t.Fatal(err)
	}
	input := make(chan Message)
	output := make(chan Message)
	
	go s.Start(input, output)
	go func() {
		time.Sleep(100 * time.Millisecond)
		gpio.On(nil)
		time.Sleep(100 * time.Millisecond)
		gpio.Off()
	}()
	val := <-output
	if val.Value.Value.(float64) != 5.0 {
		t.Error("should have been 5.0", val.Value)
	}
	val = <-output
	if val.Value.Value.(float64) != 0.0 {
		t.Error("should have been 0.0", val.Value)
	}
}

func TestInputStart(t *testing.T) {
	location := "lab"
	name := "switch"
	poller := &FakePoller{}
	s := &Switch{
		GPIO: poller,
		Value: 5.0,
		Units: "liters",
	}
	g := Gadget{
		Location: location,
		Name: name,
		Input: s,
		UID: fmt.Sprintf("%s %s", location, name),
	}
	input := make(chan Message)
	output := make(chan Message)
	go g.Start(input, output)
	val := <-output
	if val.Value.Value.(float64) != 5.0 {
		t.Error("should have been 5.0", val.Value)
	}
	val = <-output
	if val.Value.Value.(float64) != 0.0 {
		t.Error("should have been 0.0", val.Value)
	}
}
