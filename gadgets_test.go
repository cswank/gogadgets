package gogadgets

import (
	"bitbucket.org/cswank/gogadgets/utils"
	"fmt"
	"testing"
	"time"
)

func TestStripCommand(t *testing.T) {
	cmd := stripCommand("fill tank to 5 liters")
	if cmd != "5 liters" {
		t.Error(cmd)
	}
	cmd = stripCommand("turn on lab led for 2.3 minutes")
	if cmd != "2.3 minutes" {
		t.Error(cmd)
	}
}

func TestParseCommand(t *testing.T) {
	val, unit, err := ParseCommand("fill tank to 5 liters")
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
	val, unit, err := ParseCommand("turn on lab led for 1.1 minutes")
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

func TestGetDuration(t *testing.T) {
	g := Gadget{}
	d := g.getDuration(2, "minutes")
	e := time.Duration(2.0 * time.Minute)
	if d != e {
		t.Error(d, e)
	}
}

func TestGetHourDuration(t *testing.T) {
	g := Gadget{}
	d := g.getDuration(1, "hour")
	e := time.Duration(1.0 * time.Hour)
	if d != e {
		t.Error(d, e)
	}
}

func TestStart(t *testing.T) {
	location := "lab"
	name := "led"
	g := Gadget{
		Location:   location,
		Name:       name,
		Direction:  "output",
		OnCommand:  fmt.Sprintf("turn on %s %s", location, name),
		OffCommand: fmt.Sprintf("turn off %s %s", location, name),
		Output:     &FakeOutput{},
		UID:        fmt.Sprintf("%s %s", location, name),
	}
	input := make(chan Message)
	output := make(chan Message)
	go g.Start(input, output)
	update := <-output
	if update.Value.Value != false {
		t.Error("shoulda been off", update.Value.Value)
	}
	msg := Message{
		Type: "command",
		Body: "turn on lab led",
	}
	input <- msg
	update = <-output
	if update.Value.Value != true {
		t.Fatal("shoulda been on", update)
	}

	msg = Message{
		Type: "command",
		Body: "turn off lab led",
	}
	input <- msg
	update = <-output
	if update.Value.Value != false {
		t.Error("shoulda been off", update.Value.Value)
	}
	msg = Message{
		Type: "command",
		Body: "shutdown",
	}
	input <- msg
	update = <-output
}

func TestStartWithTrigger(t *testing.T) {
	location := "tank"
	name := "valve"
	g := Gadget{
		Location:   location,
		Name:       name,
		Operator:   ">=",
		OnCommand:  fmt.Sprintf("fill %s", location),
		OffCommand: fmt.Sprintf("stop filling %s", location),
		Output:     &FakeOutput{},
		UID:        fmt.Sprintf("%s %s", location, name),
	}
	input := make(chan Message)
	output := make(chan Message)
	go g.Start(input, output)
	update := <-output
	if update.Value.Value != false {
		t.Error("shoulda been off", update.Value.Value)
	}
	msg := Message{
		Type: "command",
		Body: "fill tank to 4.4 liters",
	}
	input <- msg
	update = <-output
	if update.Value.Value != true {
		t.Error("shoulda been on", update)
	}
	//make a message that should trigger the trigger and stop the device
	msg = Message{
		Sender:   "tank volume",
		Type:     UPDATE,
		Location: "tank",
		Name:     "volume",
		Value: Value{
			Units: "liters",
			Value: 4.4,
		},
	}
	input <- msg
	update = <-output
	if update.Value.Value != false {
		t.Error("shoulda been off", update)
	}
}

func TestStartWithTimeTrigger(t *testing.T) {
	location := "lab"
	name := "led"
	g := Gadget{
		Location:   location,
		Name:       name,
		OnCommand:  "turn on lab led",
		Operator:   ">=",
		OffCommand: "turn off lab led",
		Output:     &FakeOutput{},
		UID:        fmt.Sprintf("%s %s", location, name),
	}
	input := make(chan Message)
	output := make(chan Message)
	go g.Start(input, output)
	update := <-output
	if update.Value.Value != false {
		t.Error("shoulda been off", update.Value.Value)
	}
	msg := Message{
		Type: "command",
		Body: "turn on lab led for 0.01 seconds",
	}
	input <- msg
	update = <-output
	if update.Value.Value != true {
		t.Error("shoulda been on", update)
	}
	//wait for a second
	update = <-output
	if update.Value.Value != false {
		t.Error("shoulda been off", update)
	}
}

func TestStartWithTimeTriggerWithInterrupt(t *testing.T) {
	location := "lab"
	name := "led"
	g := Gadget{
		Location:   location,
		Name:       name,
		OnCommand:  "turn on lab led",
		OffCommand: "turn off lab led",
		Output:     &FakeOutput{},
		UID:        fmt.Sprintf("%s %s", location, name),
	}
	input := make(chan Message)
	output := make(chan Message)
	go g.Start(input, output)
	update := <-output
	if update.Value.Value != false {
		t.Error("shoulda been off", update.Value.Value)
	}
	msg := Message{
		Type: "command",
		Body: "turn on lab led for 30 seconds",
	}
	input <- msg
	update = <-output
	if update.Value.Value != true {
		t.Error("shoulda been on", update)
	}

	msg = Message{
		Type: "command",
		Body: "turn on lab led",
	}
	input <- msg

	msg = Message{
		Type: "update",
		Body: "",
	}
	input <- msg

	msg = Message{
		Type: "command",
		Body: "turn off lab led",
	}
	input <- msg
	update = <-output
	if update.Value.Value != false {
		t.Error("shoulda been off", update)
	}
}

func _TestStartWithTimeTriggerForReals(t *testing.T) {
	if !utils.FileExists("/sys/class/gpio/export") {
		return //not a beaglebone
	}
	pin := &Pin{Type: "gpio", Port: "9", Pin: "15"}
	gpio, err := NewGPIO(pin)
	if err != nil {
		t.Fatal(err)
	}
	location := "lab"
	name := "led"
	g := Gadget{
		Location:   location,
		Name:       name,
		OnCommand:  "turn on lab led",
		OffCommand: "turn off lab led",
		Output:     gpio,
		UID:        fmt.Sprintf("%s %s", location, name),
	}
	input := make(chan Message)
	output := make(chan Message)
	go g.Start(input, output)
	update := <-output
	if update.Value.Value != false {
		t.Error("shoulda been off", update.Value.Value)
	}
	msg := Message{
		Type: "command",
		Body: "turn on lab led for 0.1 seconds",
	}
	input <- msg
	update = <-output
	if update.Value.Value != true {
		t.Error("shoulda been on", update)
	}
	//wait for a second
	update = <-output
	if update.Value.Value != false {
		t.Error("shoulda been off", update)
	}
}

func _TestRealInput(t *testing.T) {
	if !utils.FileExists("/sys/class/gpio/export") {
		return //not a beaglebone
	}

	gpioConfig := &Pin{
		Type: "gpio",
		Port: "9",
		Pin:  "15",
	}
	gpio, err := NewGPIO(gpioConfig)
	if err != nil {
		t.Fatal(err)
	}

	config := &GadgetConfig{
		Location: "lab",
		Name:     "switch",
		Pin: Pin{
			Type:      "switch",
			Port:      "9",
			Pin:       "16",
			Edge:      "both",
			Direction: "in",
			Value:     5.0,
			Units:     "liters",
		},
	}
	s, err := NewGadget(config)
	if err != nil {
		t.Fatal(err)
	}
	input := make(chan Message)
	output := make(chan Message)

	go s.Start(input, output)
	update := <-output
	if update.Value.Value != false {
		t.Error("shoulda been off", update.Value.Value)
	}
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
		GPIO:      poller,
		Value:     5.0,
		TrueValue: 5.0,
		Units:     "liters",
	}
	g := Gadget{
		Location: location,
		Name:     name,
		Input:    s,
		UID:      fmt.Sprintf("%s %s", location, name),
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
