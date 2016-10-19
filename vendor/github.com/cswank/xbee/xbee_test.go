package xbee_test

import (
	"github.com/cswank/xbee"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestXbee(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Xbee Suite")
}

var _ = Describe("xbee message", func() {

	var (
		x xbee.Message
	)

	BeforeEach(func() {
		var err error
		d := []byte{0x92, 0x00, 0x13, 0xA2, 0x00, 0x40, 0x4C, 0x0E, 0xBE, 0x61, 0x59, 0x01, 0x01, 0x00, 0x18, 0x03, 0x00, 0x10, 0x02, 0x2F, 0x01, 0xFE, 0x49}
		x, err = xbee.NewMessage(d)
		Expect(err).To(BeNil())
	})

	It("gets the analog values", func() {
		vals, err := x.GetAnalog()
		Expect(err).To(BeNil())
		Expect(vals).To(HaveLen(2))
		Expect(vals["adc0"]).To(BeNumerically("~", 655.7, 0.1))
		Expect(vals["adc1"]).To(BeNumerically("~", 598.2, 0.1))
	})

	It("gets the digital values", func() {
		vals, err := x.GetDigital()
		Expect(err).To(BeNil())
		Expect(vals).To(HaveLen(2))
		v, ok := vals["dio3"]
		Expect(ok).To(BeTrue())
		Expect(v).To(BeFalse())
		Expect(vals["dio4"]).To(BeTrue())
	})
})
