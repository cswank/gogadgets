package gogadgets_test

import (
	"time"

	"github.com/cswank/gogadgets"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("method runner", func() {
	var (
		m       *gogadgets.MethodRunner
		in, out chan gogadgets.Message
	)

	BeforeEach(func() {
		in = make(chan gogadgets.Message)
		out = make(chan gogadgets.Message)
		m = &gogadgets.MethodRunner{}
	})

	Context("running methods", func() {
		It("runs a method with a user wait command", func() {
			go m.Start(out, in)
			msg := gogadgets.Message{
				Type: gogadgets.METHOD,
				Method: gogadgets.Method{
					Steps: []string{
						"turn on lab led",
						"wait for 0.1 seconds",
						"turn off lab led",
						"wait for user to turn off power",
						"shutdown",
					},
				},
			}
			out <- msg
			msg = <-in
			Expect(msg.Type).To(Equal("method update"))
			Expect(msg.Method.Step).To(Equal(0))

			msg = <-in
			Expect(msg.Type).To(Equal("command"))
			Expect(msg.Body).To(Equal("turn on lab led"))
			msg = <-in
			Expect(msg.Type).To(Equal("method update"))
			Expect(msg.Method.Step).To(Equal(1))
			msg = <-in
			Expect(msg.Type).To(Equal("method update"))
			Expect(msg.Method.Step).To(Equal(1))
			msg = <-in
			Expect(msg.Type).To(Equal("method update"))
			Expect(msg.Method.Step).To(Equal(1))
			<-in
			msg = <-in
			Expect(msg.Type).To(Equal("command"))
			Expect(msg.Body).To(Equal("turn off lab led"))
			msg = <-in
			Expect(msg.Type).To(Equal("method update"))
			Expect(msg.Method.Step).To(Equal(3))
			out <- gogadgets.Message{
				Type: "update",
				Body: "wait for user to turn off power",
			}
			msg = <-in
			Expect(msg.Type).To(Equal("method update"))
			Expect(msg.Method.Step).To(Equal(4))
			msg = <-in
			Expect(msg.Type).To(Equal("command"))
			Expect(msg.Body).To(Equal("shutdown"))
		})

		It("runs a method with a wait command", func() {
			go m.Start(out, in)
			msg := gogadgets.Message{
				Type: gogadgets.METHOD,
				Method: gogadgets.Method{
					Steps: []string{
						"wait for sun temperature <= 100 C",
						"shutdown",
					},
				},
			}

			out <- msg
			msg = <-in
			out <- gogadgets.Message{
				Sender: "sum temperature",
				Value: gogadgets.Value{
					Value: 99.9,
					Units: "C",
				},
			}

			var x bool
			select {
			case <-in:
				x = false
			case <-time.After(100 * time.Millisecond):
				x = true
			}
			Expect(x).To(BeTrue())

			out <- gogadgets.Message{
				Sender: "sun temperature",
				Type:   "update",
				Value: gogadgets.Value{
					Value: 100.0,
					Units: "C",
				},
			}
			msg = <-in
			Expect(msg.Type).To(Equal("method update"))
			Expect(msg.Method.Step).To(Equal(1))

			msg = <-in
			Expect(msg.Type).To(Equal("command"))
			Expect(msg.Body).To(Equal("shutdown"))
		})

		It("resumes a method with a wait for time command", func() {
			go m.Start(out, in)
			msg := gogadgets.Message{
				Type: gogadgets.METHOD,
				Method: gogadgets.Method{
					Step: 1,
					Time: 30,
					Steps: []string{
						"turn on sun",
						"wait for 60 seconds",
						"shutdown",
					},
				},
			}

			out <- msg
			msg = <-in
			Expect(msg.Type).To(Equal("method update"))
			Expect(msg.Method.Step).To(Equal(1))
			Expect(msg.Method.Time).To(Equal(30))

			msg = <-in
			Expect(msg.Type).To(Equal("method update"))
			Expect(msg.Method.Step).To(Equal(1))
			Expect(msg.Method.Time).To(Equal(30))

			msg = <-in
			Expect(msg.Type).To(Equal("method update"))
			Expect(msg.Method.Step).To(Equal(1))
			Expect(msg.Method.Time).To(Equal(29))
		})

		It("resumes a method with a wait for user command", func() {
			go m.Start(out, in)
			msg := gogadgets.Message{
				Type: gogadgets.METHOD,
				Method: gogadgets.Method{
					Step: 1,
					Steps: []string{
						"turn on sun",
						"wait for user to say ok",
						"turn off sun",
						"shutdown",
					},
				},
			}

			out <- msg
			msg = <-in
			Expect(msg.Type).To(Equal("method update"))
			Expect(msg.Method.Step).To(Equal(1))
		})
	})

	It("resumes a method", func() {
		go m.Start(out, in)
		msg := gogadgets.Message{
			Type: gogadgets.METHOD,
			Method: gogadgets.Method{
				Step: 1,
				Steps: []string{
					"turn on sun",
					"turn off sun",
					"turn on sun",
					"shutdown",
				},
			},
		}

		out <- msg
		msg = <-in
		Expect(msg.Type).To(Equal("method update"))
		Expect(msg.Method.Step).To(Equal(2))
	})
})
