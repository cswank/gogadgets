package methods

import (
	"testing"
	"time"
)

func TestRunMethod(t *testing.T) {
	// method:= []string{
	// 	"turn on lab led",
	// 	"wait for 2 seconds",
	// 	"turn off lab led",
	// }
}


func TestReadWaitCommand(t *testing.T) {
	m := Methods{}
	d, err := m.readWaitCommand("wait for 3.3 seconds")
	if err != nil {
		t.Error(err)
	}
	if d != time.Duration(3.3 * float64(time.Second)) {
		t.Error("incorrect time", d)
	}

	m = Methods{}
	d, err = m.readWaitCommand("wait for 1 second")
	if err != nil {
		t.Error(err)
	}
	if d != time.Duration(1.0 * float64(time.Second)) {
		t.Error("incorrect time", d)
	}

	m = Methods{}
	d, err = m.readWaitCommand("wait for 10 hours")
	if err != nil {
		t.Error(err)
	}
	if d != time.Duration(36000.0 * float64(time.Second)) {
		t.Error("incorrect time", d)
	}

	d, err = m.readWaitCommand("wait for 1.1 minutes")
	if err != nil {
		t.Error(err)
	}
	if d != time.Duration(1.1 * 60.0 * float64(time.Second)) {
		t.Error("incorrect time", d)
	}
}
