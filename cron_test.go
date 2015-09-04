package gogadgets_test

import (
	"math/rand"
	"time"

	"github.com/cswank/gogadgets"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type fakeNow struct {
	t time.Time
}

func (f *fakeNow) Now() time.Time {
	return t
}

func init() {
	rand.Seed(time.Now().Unix())
}

var _ = Describe("Switch", func() {
	var (
		out  chan gogadgets.Message
		in   chan gogadgets.Value
		s    *gogadgets.Switch
		fn   fakeNow
		jobs string
	)

	BeforeEach(func() {
		jobs = `
25, 13, * * * turn on living room light
`
		fn = fakeNow{
			t: time.Date(2015, 9, 4, 13, 25, 0, 0, time.UTC),
		}
		c = &gogadgets.Cron{
			Now:  fn,
			jobs: jobs,
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
