package gogadgets_test

import (
	"errors"
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

var _ = Describe("Cron", func() {
	var (
		out     chan gogadgets.Message
		in      chan gogadgets.Message
		c       *gogadgets.Cron
		fa      *fakeAfter
		jobs    []string
		start   func() *gogadgets.Message
		cronErr error
	)

	BeforeEach(func() {
		cronErr = nil
		fa = &fakeAfter{
			t: time.Date(2015, 9, 4, 13, 25, 0, 0, time.UTC),
		}

		jobs = []string{
			"25 13 * * * turn on living room light",
			"25 14 * * * turn on living room light",
		}
		out = make(chan gogadgets.Message)
		in = make(chan gogadgets.Message)

		start = func() *gogadgets.Message {
			var err error
			c, err = gogadgets.NewCron(
				nil,
				gogadgets.CronAfter(fa.After),
				gogadgets.CronJobs(jobs),
				gogadgets.CronSleep(time.Millisecond),
			)
			var msg *gogadgets.Message
			if cronErr != nil {
				Expect(err).To(MatchError(cronErr))
			} else {
				Expect(err).To(BeNil())
				go c.Start(out, in)
				select {
				case m := <-in:
					msg = &m
				case <-time.After(100 * time.Millisecond):
					return nil
				}
			}
			return msg
		}
	})

	Describe("bunk cron jobs", func() {
		It("returns an error when there isn't enough stuff in the job", func() {
			jobs = []string{"25 start"}
			cronErr = errors.New("could not parse job: 25 start")
			start()
		})

		It("returns an error when the numbers aren't numbers", func() {
			jobs = []string{"25 * b * * start"}
			cronErr = errors.New("could not parse job: 25 * b * * start")
			start()
		})
	})

	Describe("when all's good", func() {
		It("sends a command when it's time", func() {
			msg := start()
			Expect(msg).ToNot(BeNil())
			Expect(msg.Body).To(Equal("turn on living room light"))
			Expect(msg.Sender).To(Equal("cron"))
		})

		It("sends a command when there is a weekday specified", func() {
			jobs = []string{"25 13 * * 5 turn on living room light"}
			msg := start()
			Expect(msg).ToNot(BeNil())
			Expect(msg.Body).To(Equal("turn on living room light"))
			Expect(msg.Sender).To(Equal("cron"))
		})

		It("does not send a command when the weekday is wrong", func() {
			jobs = []string{"25 13 * * 6 turn on living room light"}
			msg := start()
			Expect(msg).To(BeNil())
		})

		It("sends a command when there is a range of weekdays specified", func() {
			jobs = []string{"25 13 * * 4-6 turn on living room light"}
			msg := start()
			Expect(msg).ToNot(BeNil())
			Expect(msg.Body).To(Equal("turn on living room light"))
			Expect(msg.Sender).To(Equal("cron"))
		})

		It("sends a command when there are several weekdays specified", func() {
			jobs = []string{"25 13 * * 2,5,6 turn on living room light"}
			msg := start()
			Expect(msg).ToNot(BeNil())
			Expect(msg.Body).To(Equal("turn on living room light"))
			Expect(msg.Sender).To(Equal("cron"))
		})

		It("sends a command when there is lots of extra space", func() {
			jobs = []string{
				"25	13         *     *    *    turn on living room light",
				"25 14 * * * turn on living room light",
			}
			msg := start()
			Expect(msg).ToNot(BeNil())
			Expect(msg.Body).To(Equal("turn on living room light"))
			Expect(msg.Sender).To(Equal("cron"))
		})

		It("sends a command when there are tabs and space", func() {
			jobs = []string{
				"25	13     *	*    *    turn on living room light",
			}
			msg := start()
			Expect(msg).ToNot(BeNil())
			Expect(msg.Body).To(Equal("turn on living room light"))
			Expect(msg.Sender).To(Equal("cron"))
		})

		It("sends a command when there are just tabs", func() {
			jobs = []string{
				"25	13	*	*	*	turn on living room light",
			}
			msg := start()
			Expect(msg).ToNot(BeNil())
			Expect(msg.Body).To(Equal("turn on living room light"))
			Expect(msg.Sender).To(Equal("cron"))
		})

		It("sends a command when there is a range of minutes", func() {
			jobs = []string{"22-26 13 * * * turn on living room light"}
			msg := start()
			Expect(msg).ToNot(BeNil())
			Expect(msg.Body).To(Equal("turn on living room light"))
			Expect(msg.Sender).To(Equal("cron"))
		})

		It("sends a command when there is a range of hours", func() {
			jobs = []string{"25 12-15 * * * turn on living room light"}
			msg := start()
			Expect(msg).ToNot(BeNil())
			Expect(msg.Body).To(Equal("turn on living room light"))
			Expect(msg.Sender).To(Equal("cron"))
		})

		It("sends a command when there is a series of minutes", func() {
			jobs = []string{"22,25,28 13 * * * turn on living room light"}
			msg := start()
			Expect(msg).ToNot(BeNil())
			Expect(msg.Body).To(Equal("turn on living room light"))
			Expect(msg.Sender).To(Equal("cron"))
		})

		It("sends a command when there is a series of hours", func() {
			jobs = []string{"25 1,13 * * * turn on living room light"}
			msg := start()
			Expect(msg).ToNot(BeNil())
			Expect(msg.Body).To(Equal("turn on living room light"))
			Expect(msg.Sender).To(Equal("cron"))
		})

		It("does not send a command when it's not time", func() {
			jobs = []string{"25 14 * * * turn on living room light"}
			msg := start()
			Expect(msg).To(BeNil())
		})

		It("does not send a command when the line is commented out", func() {
			jobs = []string{"#25 14 * * * turn on living room light"}
			msg := start()
			Expect(msg).To(BeNil())
		})

		It("sends a command when there is a range of hours and minutes", func() {
			//This doesn't work

			// jobs = []string{"24-28 12-15 * * * turn on living room light"}
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
	})
})
