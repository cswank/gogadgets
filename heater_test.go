package gogadgets

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"os"
// 	"testing"
// 	"time"

// 	"github.com/cswank/gogadgets/utils"
// )

// var (
// 	testGPIODevPath = "/tmp/sys/class/gpio/gpio45"
// )

// func init() {
// 	if !utils.FileExists(testGPIODevPath) {
// 		os.MkdirAll(testGPIODevPath, 0777)
// 	}
// }

// func getMessage(val float64) *Message {
// 	return &Message{
// 		Name: "temperature",
// 		Value: Value{
// 			Value: val,
// 		},
// 	}
// }

// func getValue(pth string) string {
// 	d, _ := ioutil.ReadFile(fmt.Sprintf("%s/%s", testGPIODevPath, pth))
// 	return string(d)
// }

// func waitFor(f, val string) {
// 	v := getValue(f)
// 	for v != val {
// 		time.Sleep(10 * time.Millisecond)
// 		v = getValue(f)
// 	}
// }

// func _TestPWMHeater(t *testing.T) {
// 	pwmMode = 0777
// 	GPIO_DEV_PATH = "/tmp/sys/class/gpio"
// 	GPIO_DEV_MODE = 0777
// 	p := &Pin{
// 		Type:      "heater",
// 		Port:      "8",
// 		Pin:       "11",
// 		Frequency: 1,
// 		Args:      map[string]interface{}{"pwm": true},
// 	}
// 	d, err := NewHeater(p)
// 	if err != nil {
// 		t.Fatal(err, d)
// 	}
// 	d.On(nil)
// 	waitFor("value", "1")
// 	d.Off()
// 	waitFor("value", "0")
// 	m := &Message{
// 		Name:  "temperature",
// 		Value: Value{Value: 29.0, Units: "C"},
// 	}
// 	d.Update(m)
// 	d.On(&Value{Value: 30.0, Units: "C"})
// 	time.Sleep(10 * time.Second)
// 	// waitFor("run", "1")
// 	// m = &Message{
// 	// 	Name:  "temperature",
// 	// 	Value: Value{Value: 30.0, Units: "C"},
// 	// }
// 	// waitFor("duty", "1000000000")
// 	// d.Update(m)
// 	// waitFor("duty", "0")
// 	// m = &Message{
// 	// 	Name:  "temperature",
// 	// 	Value: Value{Value: 29.0, Units: "C"},
// 	// }
// 	// d.Update(m)
// 	// waitFor("duty", "250000000")
// 	// waitFor("run", "1")
// }

// // func TestHeater(t *testing.T) {
// // 	pwmMode = 0777
// // 	GPIO_DEVPATH = "/tmp/sys/class/gpio"
// // 	p := &Pin{
// // 		Type:      "heater",
// // 		Port:      "8",
// // 		Pin:       "13",
// // 		Frequency: 1,
// // 	}
// // 	d, err := NewHeater(p)
// // 	if err != nil {
// // 		t.Error(err, d)
// // 	}
// // 	d.On(nil)
// // 	waitFor("run", "1")
// // 	d.Off()
// // 	waitFor("run", "0")
// // 	m := &Message{
// // 		Name:  "temperature",
// // 		Value: Value{Value: 20.0, Units: "C"},
// // 	}
// // 	d.Update(m)
// // 	d.On(&Value{Value: 30.0, Units: "C"})
// // 	waitFor("run", "1")
// // 	m = &Message{
// // 		Name:  "temperature",
// // 		Value: Value{Value: 30.0, Units: "C"},
// // 	}
// // 	waitFor("duty", "1000000000")
// // 	d.Update(m)
// // 	waitFor("run", "0")
// // 	m = &Message{
// // 		Name:  "temperature",
// // 		Value: Value{Value: 29.0, Units: "C"},
// // 	}
// // 	d.Update(m)
// // 	waitFor("duty", "1000000000")
// // 	waitFor("run", "1")
// // }
