package gogadgets_test

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net/http"
	"net/http/httptest"
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

	Describe("with a moisture guard", func() {
		var (
			srv   *httptest.Server
			value float64
		)

		BeforeEach(func() {
			srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				msgs := []gogadgets.Message{{
					Location: "garden bed",
					Name:     "soil moisture",
					Value:    gogadgets.Value{Value: value, Units: "%"},
				}}
				json.NewEncoder(w).Encode(msgs)
			}))
			jobs = []string{"25 13 * * * turn on garden bed sprinklers"}
		})

		AfterEach(func() {
			srv.Close()
		})

		runWithGuard := func(guard gogadgets.Guard) *gogadgets.Message {
			var err error
			c, err = gogadgets.NewCron(
				nil,
				gogadgets.CronAfter(fa.After),
				gogadgets.CronJobs(jobs),
				gogadgets.CronSleep(time.Millisecond),
				gogadgets.CronGuards([]gogadgets.Guard{guard}),
			)
			Expect(err).To(BeNil())
			go c.Start(out, in)
			select {
			case m := <-in:
				return &m
			case <-time.After(100 * time.Millisecond):
				return nil
			}
		}

		It("skips the command when sensor reads at or above max", func() {
			value = 75
			msg := runWithGuard(gogadgets.Guard{
				Match:  "turn on garden bed sprinklers",
				URL:    srv.URL,
				Device: "garden bed soil moisture",
				Max:    60,
			})
			Expect(msg).To(BeNil())
		})

		It("emits the command when sensor reads below max", func() {
			value = 30
			msg := runWithGuard(gogadgets.Guard{
				Match:  "turn on garden bed sprinklers",
				URL:    srv.URL,
				Device: "garden bed soil moisture",
				Max:    60,
			})
			Expect(msg).ToNot(BeNil())
			Expect(msg.Body).To(Equal("turn on garden bed sprinklers"))
		})

		It("fails open when the sensor fetch errors", func() {
			srv.Close()
			value = 999
			msg := runWithGuard(gogadgets.Guard{
				Match:  "turn on garden bed sprinklers",
				URL:    srv.URL,
				Device: "garden bed soil moisture",
				Max:    60,
			})
			Expect(msg).ToNot(BeNil())
			Expect(msg.Body).To(Equal("turn on garden bed sprinklers"))
		})

		It("does not gate commands whose prefix doesn't match", func() {
			jobs = []string{"25 13 * * * turn on front yard sprinklers"}
			value = 999
			msg := runWithGuard(gogadgets.Guard{
				Match:  "turn on garden bed sprinklers",
				URL:    srv.URL,
				Device: "garden bed soil moisture",
				Max:    60,
			})
			Expect(msg).ToNot(BeNil())
			Expect(msg.Body).To(Equal("turn on front yard sprinklers"))
		})

		It("matches commands with duration suffixes", func() {
			jobs = []string{"25 13 * * * turn on garden bed sprinklers for 10 minutes"}
			value = 75
			msg := runWithGuard(gogadgets.Guard{
				Match:  "turn on garden bed sprinklers",
				URL:    srv.URL,
				Device: "garden bed soil moisture",
				Max:    60,
			})
			Expect(msg).To(BeNil())
		})
	})
})
