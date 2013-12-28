package gogadgets

import (
	"fmt"
	"testing"
	"bitbucket.com/cswank/gogadgets/utils"
)

func TestThermometer(t *testing.T) {
	if !utils.FileExists("/sys/class/gpio/export") {
		return //not a beaglebone
	}
	pin := &Pin{
		OneWireId: "28-0000047ade8f",
	}
	therm, err := NewThermometer(pin)
	if err != nil {
		t.Error(err)
	}
	out := make(chan Message)
	in := make(chan Value)
	go therm.Start(out, in)
	val := <-in
	if val.Units != "C" {
		t.Error("units should have been 'C'", val)
	}
	fmt.Println("the temperature is:", val.Value)
	out<- Message{
		Type: "command",
		Body: "shutdown",
	}
}

func TestThermometerParseValue(t *testing.T) {
	therm := Thermometer{}
	val, err := therm.parseValue("3d 01 4b 46 7f ff 03 10 6d : crc=6d YES\n3d 01 4b 46 7f ff 03 10 6d t=19812\n")
	if err != nil {
		t.Error(err)
	}
	if val.Value != 19.812 {
		t.Error("incorrect val", val)
	}
}
