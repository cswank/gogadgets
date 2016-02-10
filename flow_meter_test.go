package gogadgets_test

import (
	"time"

	"github.com/cswank/gogadgets"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("flow meter", func() {
	var (
		out     chan gogadgets.Message
		in      chan gogadgets.Value
		m       *gogadgets.FlowMeter
		trigger chan bool
	)
	BeforeEach(func() {
		trigger = make(chan bool)
		poller := &FakePoller{trigger: trigger}
		m = &gogadgets.FlowMeter{
			GPIO:    poller,
			Value:   5.0,
			Units:   "liters/minute",
			MinSpan: 0.001,
		}
		out = make(chan gogadgets.Message)
		in = make(chan gogadgets.Value)
		go m.Start(out, in)
	})

	Describe("normal operation", func() {
		It("measures flow", func() {
			val := <-in
			Expect(val.Value.(float64)).To(Equal(0.0))
			trigger <- true
			time.Sleep(100 * time.Millisecond)
			trigger <- true
			val = <-in
			Expect(val.Value.(float64)).Should(BeNumerically("~", 50.0, 2.0))
			time.Sleep(100 * time.Millisecond)
			trigger <- true
			val = <-in
			Expect(val.Value.(float64)).Should(BeNumerically("~", 50.0, 2.0))
		})

	})
})
