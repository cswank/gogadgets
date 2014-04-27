package gogadgets

import (
	"bitbucket.org/cswank/gogadgets/utils"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

var (
	testDevPath = "/tmp/sys/devices/ocp.3/pwm_test_P8_13.11"
)

func init() {
	if !utils.FileExists(testDevPath) {
		os.MkdirAll(testDevPath, 0777)
	}
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
	d, _ := ioutil.ReadFile(fmt.Sprintf("%s/%s", testDevPath, pth))
	return string(d)
}

func waitFor(f, val string) {
	v := getValue(f)
	for v != val {
		v = getValue(f)
		time.Sleep(10 * time.Millisecond)
	}
}

func TestHeater(t *testing.T) {
	pwmMode = 0777
	PWM_DEVPATH = "/tmp/sys/devices/ocp.*/pwm_test_P%s_%s.*"
	p := &Pin{
		Type:      "heater",
		Port:      "8",
		Pin:       "13",
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
		Name:  "temperature",
		Value: Value{Value: 20.0, Units: "C"},
	}
	d.Update(m)
	d.On(&Value{Value: 30.0, Units: "C"})
	waitFor("run", "1")
	m = &Message{
		Name:  "temperature",
		Value: Value{Value: 30.0, Units: "C"},
	}
	waitFor("duty", "1000000000")
	d.Update(m)
	waitFor("run", "0")
	m = &Message{
		Name:  "temperature",
		Value: Value{Value: 29.0, Units: "C"},
	}
	d.Update(m)
	waitFor("duty", "1000000000")
	waitFor("run", "1")
}

func TestPWMHeater(t *testing.T) {
	pwmMode = 0777
	PWM_DEVPATH = "/tmp/sys/devices/ocp.*/pwm_test_P%s_%s.*"
	p := &Pin{
		Type:      "heater",
		Port:      "8",
		Pin:       "13",
		Frequency: 1,
		Args:      map[string]interface{}{"pwm": true},
	}
	d, err := NewHeater(p)
	if err != nil {
		t.Fatal(err, d)
	}
	d.On(nil)
	waitFor("run", "1")
	d.Off()
	waitFor("run", "0")
	m := &Message{
		Name:  "temperature",
		Value: Value{Value: 20.0, Units: "C"},
	}
	d.Update(m)
	d.On(&Value{Value: 30.0, Units: "C"})
	waitFor("run", "1")
	m = &Message{
		Name:  "temperature",
		Value: Value{Value: 30.0, Units: "C"},
	}
	waitFor("duty", "1000000000")
	d.Update(m)
	waitFor("duty", "0")
	m = &Message{
		Name:  "temperature",
		Value: Value{Value: 29.0, Units: "C"},
	}
	d.Update(m)
	waitFor("duty", "250000000")
	waitFor("run", "1")
}
