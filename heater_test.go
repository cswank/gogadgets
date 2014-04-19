package gogadgets

import (
	"os"
	"fmt"
	"time"
	"io/ioutil"
	"testing"
)

func init() {
	os.Setenv("PWMTESTDEVICEMODE", "perm")
}

func getMessage(val float64) *Message {
	return &Message{
		Name: "temperature",
		Value: Value{
			Value: val,
		},
	}
}

func getValue(pth string) string {
	d, _ := ioutil.ReadFile(fmt.Sprintf("/sys/devices/ocp.9/pwm_test_8_13.x/%s", pth))
	return string(d)
}

func waitFor(f, val string) {
	v := getValue(f)
	fmt.Println("wait for", f, val, v)
	for v != val {
		v = getValue(f)
		time.Sleep(10 * time.Millisecond)
	}
}

func _TestHeater(t *testing.T) {
	p := &Pin{
		Type: "heater",
		Port: "8",
		Pin: "13",
		Frequency: 1,
	}
	d, err := NewHeater(p)
	if err != nil {
		t.Error(err, d)
	}
	d.On(nil)
	waitFor("run", "1")
	d.Off()
	waitFor("run", "0")
	m := &Message{
		Name: "temperature",
		Value: Value{Value: 20.0, Units: "C"},
	}
	d.Update(m)
	d.On(&Value{Value:30.0, Units: "C"})
	waitFor("run", "1")
	m = &Message{
		Name: "temperature",
		Value: Value{Value: 30.0, Units: "C"},
	}
	waitFor("duty", "1000000000")
	d.Update(m)
	waitFor("run", "0")
	m = &Message{
		Name: "temperature",
		Value: Value{Value: 29.0, Units: "C"},
	}
	d.Update(m)
	waitFor("duty", "1000000000")
	waitFor("run", "1")
}

func TestPWMHeater(t *testing.T) {
	p := &Pin{
		Type: "heater",
		Port: "8",
		Pin: "13",
		Frequency: 1,
		Args: map[string]string{"pwm":"true"},
	}
	d, err := NewHeater(p)
	if err != nil {
		t.Error(err, d)
	}
	d.On(nil)
	waitFor("run", "1")
	d.Off()
	waitFor("run", "0")
	m := &Message{
		Name: "temperature",
		Value: Value{Value: 20.0, Units: "C"},
	}
	d.Update(m)
	d.On(&Value{Value:30.0, Units: "C"})
	waitFor("run", "1")
	m = &Message{
		Name: "temperature",
		Value: Value{Value: 30.0, Units: "C"},
	}
	waitFor("duty", "1000000000")
	d.Update(m)
	waitFor("run", "0")
	m = &Message{
		Name: "temperature",
		Value: Value{Value: 29.0, Units: "C"},
	}
	d.Update(m)
	waitFor("duty", "250000000")
	waitFor("run", "1")
}

