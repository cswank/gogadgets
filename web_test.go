package gogadgets_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/cswank/gogadgets"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var _ = Describe("server", func() {
	var (
		port    int
		addr    string
		cliAddr string
		out     chan gogadgets.Message
		in      chan gogadgets.Message
		s       *gogadgets.Server
		lg      *fakeLogger
		sink    []gogadgets.Message
	)

	BeforeEach(func() {
		sink = []gogadgets.Message{}
		lg = &fakeLogger{}

		port = 1024 + rand.Intn(65535-1024)
		addr = fmt.Sprintf("http://localhost:%d/gadgets", port)
		cliAddr = fmt.Sprintf("http://localhost:%d/clients", port)

		s = gogadgets.NewServer("", "", port, lg)

		in = make(chan gogadgets.Message)
		out = make(chan gogadgets.Message)
		go s.Start(out, in)
		out <- gogadgets.Message{
			Type:     gogadgets.UPDATE,
			Sender:   "lab led",
			Location: "lab",
			Name:     "led",
			Value: gogadgets.Value{
				Value:  true,
				Output: true,
			},
		}
		out <- gogadgets.Message{
			Type:     gogadgets.UPDATE,
			Sender:   "hall led",
			Location: "hall",
			Name:     "led",
			Value: gogadgets.Value{
				Value:  false,
				Output: false,
			},
		}
	})
	Describe("when all's good", func() {
		It("sends the status", func() {
			r, err := http.Get(addr)

			Expect(err).To(BeNil())
			defer r.Body.Close()

			Expect(r.StatusCode).To(Equal(http.StatusOK))
			msgs := map[string]gogadgets.Message{}
			dec := json.NewDecoder(r.Body)
			err = dec.Decode(&msgs)
			Expect(err).To(BeNil())
			Expect(len(msgs)).To(Equal(2))
			msg, ok := msgs["lab led"]
			Expect(ok).To(BeTrue())
			Expect(msg.Value.Value).To(BeTrue())

			msg, ok = msgs["hall led"]
			Expect(ok).To(BeTrue())
			Expect(msg.Value.Value).To(BeFalse())
		})
		It("accepts a message from the outside world", func() {
			msg := gogadgets.Message{
				Type:   gogadgets.COMMAND,
				Sender: "me",
				Body:   "turn on lab led",
			}

			buf := &bytes.Buffer{}
			enc := json.NewEncoder(buf)
			err := enc.Encode(&msg)
			Expect(err).To(BeNil())

			Eventually(func() int {
				r, err := http.Post(addr, "application/json", buf)
				if err != nil {
					return 500
				}
				r.Body.Close()
				return r.StatusCode
			}).Should(Equal(http.StatusOK))
			m := <-in
			Expect(m.Body).To(Equal("turn on lab led"))
		})
		It("registers a new client", func() {
			msgs := []gogadgets.Message{}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				var msg gogadgets.Message
				dec := json.NewDecoder(r.Body)
				dec.Decode(&msg)
				msgs = append(msgs, msg)
			}))
			defer ts.Close()

			a := map[string]string{"address": ts.URL}
			buf := &bytes.Buffer{}
			enc := json.NewEncoder(buf)
			enc.Encode(&a)

			Eventually(func() int {
				r, err := http.Post(cliAddr, "application/json", buf)
				if err != nil {
					return 500
				}
				r.Body.Close()
				return r.StatusCode
			}).Should(Equal(http.StatusOK))

			r, err := http.Get(cliAddr)
			Expect(err).To(BeNil())
			Expect(r.StatusCode).To(Equal(http.StatusOK))
			var c map[string]bool
			dec := json.NewDecoder(r.Body)
			err = dec.Decode(&c)
			Expect(err).To(BeNil())
			r.Body.Close()

			Expect(c[ts.URL]).To(BeTrue())

			msg := gogadgets.Message{
				Type:   gogadgets.COMMAND,
				Sender: "me",
				Body:   "turn on lab led",
			}
			out <- msg
			Eventually(func() []gogadgets.Message {
				return msgs
			}).Should(HaveLen(1))
			Expect(msgs[0].Body).To(Equal("turn on lab led"))
		})
	})
})
