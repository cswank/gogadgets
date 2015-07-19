package gogadgets

// import (
// 	"testing"
// 	"time"
// )

// func TestCreate(t *testing.T) {
// 	pin := &Pin{
// 		Args: map[string]interface{}{
// 			"host": "localhost",
// 			"db":   "brewery",
// 			"summarize": map[string]interface{}{
// 				"greenhouse temperature": 10,
// 			},
// 		},
// 	}
// 	r, err := NewRecorder(pin)
// 	if err != nil {
// 		t.Error(err, r)
// 	}
// }

// func TestSummarize(t *testing.T) {
// 	r := Recorder{
// 		summaries: map[string]time.Duration{"greenhouse temperature": time.Minute},
// 		history:   map[string]summary{},
// 	}
// 	for i := 1; i < 11; i++ {
// 		msg := &Message{
// 			Sender: "me",
// 			Value:  Value{Value: float64(i)},
// 		}
// 		r.summarize(msg, time.Minute)
// 	}
// 	if r.history["me"].v != 55.0 {
// 		t.Error(r.history["me"].v)
// 	}
// }

// func TestGetFilter(t *testing.T) {
// 	s := []string{"lab led", "lab temperature"}
// 	f := getFilter(s)
// 	if len(f) != 2 {
// 		t.Error(f)
// 	}
// 	if f[0] != "lab led" {
// 		t.Error(f)
// 	}
// 	s = []string{}
// 	f = getFilter(s)
// 	if len(f) != 0 {
// 		t.Error(f)
// 	}
// }
