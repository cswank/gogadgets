package gogadgets

// import (
// 	//"github.com/cswank/gogadgets/utils"
// 	"testing"
// 	//"time"
// )

// func TestPWM(t *testing.T) {
// 	// if !utils.FileExists("/sys/class/gpio/export") {
// 	// 	return //not a beaglebone
// 	// }
// 	// p := &Pin{
// 	// 	Port:      "8",
// 	// 	Pin:       "13",
// 	// 	Frequency: 1,
// 	// }
// 	// pwm, err := NewPWM(p)
// 	// if err != nil {
// 	// 	t.Fatal(err)
// 	// }
// 	// err = pwm.On(&Value{Value: 50, Units: "%"})
// 	// if err != nil {
// 	// 	t.Error(err)
// 	// }
// 	// time.Sleep(5 * time.Second)
// 	// err = pwm.On(&Value{Value: 10, Units: "%"})
// 	// if err != nil {
// 	// 	t.Error(err)
// 	// }
// 	// time.Sleep(5 * time.Second)
// 	// err = pwm.On(&Value{Value: 90, Units: "%"})
// 	// if err != nil {
// 	// 	t.Error(err)
// 	// }
// 	// time.Sleep(5 * time.Second)

// }

// func TestGetDuty(t *testing.T) {
// 	pwm := PWM{
// 		period: 1000000000,
// 	}
// 	duty := pwm.getDuty(50.0)
// 	if string(duty) != "500000000" {
// 		t.Error(string(duty))
// 	}
// }
