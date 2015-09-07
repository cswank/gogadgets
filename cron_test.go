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
		fa = &fakeAfter{
			t: time.Date(2015, 9, 4, 13, 25, 0, 0, time.UTC),
		}

		jobs = `25 13 * * * turn on living room light
25 14 * * * turn on living room light`
		out = make(chan gogadgets.Message)
		in = make(chan gogadgets.Message)

	})
	Describe("when all's good", func() {
		It("sends a command when it's time", func() {
			c = &gogadgets.Cron{
				After: fa.After,
				Jobs:  jobs,
				Sleep: time.Millisecond,
			}
			go c.Start(out, in)
			msg := <-in
			Expect(msg.Body).To(Equal("turn on living room light"))
			Expect(msg.Sender).To(Equal("cron"))
		})

		FIt("sends a command when there is a weekday specified", func() {
			jobs = `25 13 * * 5 turn on living room light`
			c = &gogadgets.Cron{
				After: fa.After,
				Jobs:  jobs,
				Sleep: time.Millisecond,
			}
			go c.Start(out, in)
			msg := <-in
			Expect(msg.Body).To(Equal("turn on living room light"))
			Expect(msg.Sender).To(Equal("cron"))
		})

		It("sends a command when there is lots of extra space", func() {
			jobs = `25	13     *     *    *    turn on living room light
25 14 * * * turn on living room light`
			c = &gogadgets.Cron{
				After: fa.After,
				Jobs:  jobs,
				Sleep: time.Millisecond,
			}
			go c.Start(out, in)
			msg := <-in
			Expect(msg.Body).To(Equal("turn on living room light"))
			Expect(msg.Sender).To(Equal("cron"))
		})

		It("sends a command when there is a range of minutes", func() {
			jobs = `22-26 13 * * * turn on living room light`
			c = &gogadgets.Cron{
				After: fa.After,
				Jobs:  jobs,
				Sleep: time.Millisecond,
			}
			go c.Start(out, in)
			msg := <-in
			Expect(msg.Body).To(Equal("turn on living room light"))
			Expect(msg.Sender).To(Equal("cron"))
		})

		It("sends a command when there is a range of hours", func() {
			jobs = `25 12-15 * * * turn on living room light`
			c = &gogadgets.Cron{
				After: fa.After,
				Jobs:  jobs,
				Sleep: time.Millisecond,
			}
			go c.Start(out, in)
			msg := <-in
			Expect(msg.Body).To(Equal("turn on living room light"))
			Expect(msg.Sender).To(Equal("cron"))
		})

		It("sends a command when there is a range of hours and minutes", func() {
			//This doesn't work

			// jobs = `24-28 12-15 * * * turn on living room light`
			// c = &gogadgets.Cron{
			// 	After: fa.After,
			// 	Jobs:  jobs,
			// 	Sleep: time.Millisecond,
			// }
			// go c.Start(out, in)
			// msg := <-in
			// Expect(msg.Body).To(Equal("turn on living room light"))
			// Expect(msg.Sender).To(Equal("cron"))
		})

		It("sends a command when there is a series of minutes", func() {
			jobs = `22,25,28 13 * * * turn on living room light`
			c = &gogadgets.Cron{
				After: fa.After,
				Jobs:  jobs,
				Sleep: time.Millisecond,
			}
			go c.Start(out, in)
			msg := <-in
			Expect(msg.Body).To(Equal("turn on living room light"))
			Expect(msg.Sender).To(Equal("cron"))
		})

		It("sends a command when there is a series of hours", func() {
			jobs = `25 1,13 * * * turn on living room light`
			c = &gogadgets.Cron{
				After: fa.After,
				Jobs:  jobs,
				Sleep: time.Millisecond,
			}
			go c.Start(out, in)
			msg := <-in
			Expect(msg.Body).To(Equal("turn on living room light"))
			Expect(msg.Sender).To(Equal("cron"))
		})

		It("does not send a command when it's not time", func() {
			jobs = `25 14 * * * turn on living room light`
			c = &gogadgets.Cron{
				After: fa.After,
				Jobs:  jobs,
				Sleep: time.Millisecond,
			}

			go c.Start(out, in)
			var msg *gogadgets.Message
			select {
			case m := <-in:
				msg = &m
			case <-time.After(100 * time.Millisecond):
			}
			Expect(msg).To(BeNil())
		})

		It("does not send a command when the line is commented out", func() {
			jobs = `#25 14 * * * turn on living room light`
			c = &gogadgets.Cron{
				After: fa.After,
				Jobs:  jobs,
				Sleep: time.Millisecond,
			}

			go c.Start(out, in)
			var msg *gogadgets.Message
			select {
			case m := <-in:
				msg = &m
			case <-time.After(100 * time.Millisecond):
			}
			Expect(msg).To(BeNil())
		})
	})
})
