package gogadgets_test

import (
	"math/rand"
	"time"

	"github.com/cswank/gogadgets"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type fakeAfter struct {
	t time.Time
	d time.Duration
	c chan time.Time
}

func (f *fakeAfter) After(d time.Duration) <-chan time.Time {
	f.d = d
	f.c = make(chan time.Time)
	go f.start()
	return f.c
}

func (f *fakeAfter) start() {
	time.Sleep(f.d)
	f.c <- f.t
}

func init() {
	rand.Seed(time.Now().Unix())
}

var _ = Describe("Switch", func() {
	var (
		out  chan gogadgets.Message
		in   chan gogadgets.Message
		c    *gogadgets.Cron
		fa   *fakeAfter
		jobs string
	)

	BeforeEach(func() {
		jobs = `
25 13 * * * turn on living room light
`
		fa = &fakeAfter{
			t: time.Date(2015, 9, 4, 13, 25, 0, 0, time.UTC),
		}

		c = &gogadgets.Cron{
			After: fa.After,
			Jobs:  jobs,
			Sleep: time.Millisecond,
		}
		out = make(chan gogadgets.Message)
		in = make(chan gogadgets.Message)
		go c.Start(out, in)
	})
	Describe("when all's good", func() {
		It("sends a command when its time", func() {
			msg := <-in
			Expect(msg.Body).To(Equal("turn on living room light"))
			Expect(msg.Sender).To(Equal("cron"))
		})
	})
})
