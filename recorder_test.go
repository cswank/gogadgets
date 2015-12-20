package gogadgets

import (
	"github.com/cswank/gogadgets"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Recorder", func() {
	var (
		out chan gogadgets.Message
		in  chan gogadgets.Value
		r   gogadgets.OutputDevice
	)

	BeforeEach(func() {
		cfg := `{
    "location": "lab",
    "name": "recorder",
    "initialValue": "turn on lab recorder",
    "pin": {
      "type": "recorder",
      "args": {
        "host": "localhost",
        "summarize": "1"
      }
   }
}`
		p := &gogadgets.Pin{}
		r = gogadgets.NewRecorder(p)
		out = make(chan gogadgets.Message)
		in = make(chan gogadgets.Value)
		go r.Start(out, in)
	})
	Describe("when all's good", func() {
		It("does it's thing", func() {
			Expect(1).To(Equal(1))
		})
	})
})
