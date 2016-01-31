package gogadgets_test

import (
	"github.com/cswank/gogadgets"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tank", func() {
	var (
		pin *gogadgets.Pin
		t   gogadgets.InputDevice
		out chan gogadgets.Message
		in  chan gogadgets.Value
	)
	BeforeEach(func() {
		out = make(chan gogadgets.Message)
		in = make(chan gogadgets.Value)
	})
	AfterEach(func() {

	})
	Context("valve source", func() {
		BeforeEach(func() {
			pin = &gogadgets.Pin{
				Args: map[string]interface{}{
					"type":         "valve",
					"float_switch": "hlt float switch",
					"full":         8.0,
				},
			}
			var err error
			t, err = gogadgets.NewTank(pin)
			Expect(err).To(BeNil())
		})
		It("creates a tank", func() {
			Expect(t).ToNot(BeNil())
		})
		It("responds to the float switch message to fill itself", func() {
			go t.Start(out, in)
			out <- gogadgets.Message{
				Sender: "hlt float switch",
				Value: gogadgets.Value{
					Value: true,
				},
			}
			v := <-in
			Expect(v.Value).To(Equal(8.0))
		})
	})
	Context("tank source", func() {
		BeforeEach(func() {
			pin = &gogadgets.Pin{
				Args: map[string]interface{}{
					"type":   "tank",
					"source": "hlt tank",
				},
			}
			var err error
			t, err = gogadgets.NewTank(pin)
			Expect(err).To(BeNil())
		})
		It("creates a tank", func() {
			Expect(t).ToNot(BeNil())
		})
		It("responds to hlt tank messages to fill itself", func() {
			go t.Start(out, in)
			for i := 1; i <= 8; i++ {
				out <- gogadgets.Message{
					Sender: "hlt tank",
					Value: gogadgets.Value{
						Diff: -0.5,
					},
				}
				v := <-in
				Expect(v.Value).To(Equal(-1 * 0.5 * float64(i)))
			}
		})
	})
})
