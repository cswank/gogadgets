package gogadgets_test

import (
	"math/rand"
	"time"

	"github.com/cswank/gogadgets"
	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"
)

func init() {
	rand.Seed(time.Now().Unix())
}

var _ = Describe("Switch", func() {
	var (
		out chan gogadgets.Message
		in  chan gogadgets.Value
		s   *gogadgets.Switch
	)

	BeforeEach(func() {
		poller := &FakePoller{}
		s = &gogadgets.Switch{
			GPIO: poller,
		}
		out = make(chan gogadgets.Message)
		in = make(chan gogadgets.Value)
		go s.Start(out, in)
	})
	Describe("when all's good", func() {
		It("does it's thing", func() {

		})
	})
})

// func TestBoolSwitch(t *testing.T) {
// 	poller := &FakePoller{}
// 	s := &Switch{
// 		GPIO:      poller,
// 		Value:     true,
// 		TrueValue: true,
// 	}
// 	out := make(chan Message)
// 	in := make(chan Value)
// 	go s.Start(out, in)
// 	val := <-in
// 	if val.Value != true {
// 		t.Error("should have been true", val)
// 	}
// 	val = <-in
// 	if val.Value != false {
// 		t.Error("should have been false", val)
// 	}
// 	out <- Message{
// 		Type: "command",
// 		Body: "shutdown",
// 	}
// 	v := s.GetValue()
// 	if v.Value != false {
// 		t.Error("should have been false", v)
// 	}
// }
