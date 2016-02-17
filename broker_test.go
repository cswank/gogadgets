package gogadgets_test

import (
	"fmt"
	"sync"

	"github.com/cswank/gogadgets"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type fakeGadget struct {
	in   chan gogadgets.Message
	msgs []gogadgets.Message
	lock sync.Mutex
}

func (f *fakeGadget) add(msg gogadgets.Message) {
	f.lock.Lock()
	f.msgs = append(f.msgs, msg)
	f.lock.Unlock()
}

func (f *fakeGadget) len() int {
	f.lock.Lock()
	i := len(f.msgs)
	f.lock.Unlock()
	return i
}

func (f *fakeGadget) start() {
	for {
		msg := <-f.in
		f.add(msg)
	}
}

var _ = Describe("broker", func() {
	var (
		b       *gogadgets.Broker
		clients map[string]chan gogadgets.Message
		in, out chan gogadgets.Message
		gadgets []*fakeGadget
	)

	BeforeEach(func() {
		in = make(chan gogadgets.Message)
		out = make(chan gogadgets.Message)
		clients = map[string]chan gogadgets.Message{
			"x": make(chan gogadgets.Message),
			"y": make(chan gogadgets.Message),
		}

		gadgets = []*fakeGadget{}
		for _, n := range []string{"x", "y"} {
			g := &fakeGadget{in: clients[n]}
			go g.start()
			gadgets = append(gadgets, g)
		}
		b = gogadgets.NewBroker(clients, in, out)
		go b.Start()

	})

	Context("brokering", func() {
		It("brokers", func() {
			for i := 0; i < 10; i++ {
				out <- gogadgets.Message{
					Sender: "the test",
					Type:   "junk",
					Body:   fmt.Sprintf("%d", i),
				}
			}
			Eventually(func() bool {
				return gadgets[0].len() == 10 && gadgets[1].len() == 10
			}).Should(BeTrue())
		})
	})
})
