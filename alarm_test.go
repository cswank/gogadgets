package gogadgets_test

import (
	"io/ioutil"
	"math/rand"
	"time"

	"github.com/cswank/gogadgets"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func init() {
	rand.Seed(time.Now().Unix())
}

var _ = Describe("Alarm", func() {
	var (
		a         *gogadgets.Alarm
		tmp       string
		sys       map[string]string
		status    bool
		frontDoor bool
		delay     string
	)

	BeforeEach(func() {
		var err error
		delay = "10ms"
		tmp, err = ioutil.TempDir("", "")
		gogadgets.GPIO_DEV_PATH = tmp
		gogadgets.GPIO_DEV_MODE = 0777

		Expect(err).To(BeNil())
		m := map[string]string{
			"alarm": gogadgets.Pins["gpio"]["8"]["11"],
		}
		sys = setupGPIOs(tmp, m)
	})

	JustBeforeEach(func() {
		p := &gogadgets.Pin{
			Type: "alarm",
			Args: map[string]interface{}{
				"events": map[string]interface{}{
					"front door": false,
				},
				"duration": "100ms",
				"delay":    delay,
			},
			Pins: map[string]gogadgets.Pin{
				"gpio": {
					Type: "gpio",
					Port: "8",
					Pin:  "11",
				},
			},
		}

		o, err := gogadgets.NewOutputDevice(p)
		Expect(err).To(BeNil())

		var ok bool
		a, ok = o.(*gogadgets.Alarm)
		Expect(ok).To(BeTrue())

		if status {
			a.On(nil)
		}

		msg := &gogadgets.Message{
			Sender: "front door",
			Value: gogadgets.Value{
				Value: frontDoor,
			},
		}

		a.Update(msg)
	})

	Context("armed", func() {
		BeforeEach(func() {
			frontDoor = false
			status = true
		})

		It("turns on the gpio when the front door is open", func() {
			Eventually(func() string {
				b, _ := ioutil.ReadFile(sys["alarm-value"])
				return string(b)
			}).Should(Equal("1"))

			Eventually(func() string {
				b, _ := ioutil.ReadFile(sys["alarm-value"])
				return string(b)
			}).Should(Equal("0"))
		})
	})

	Context("not armed", func() {
		BeforeEach(func() {
			frontDoor = false
			status = false
		})

		It("turns on the gpio when the front door is open", func() {
			b, err := ioutil.ReadFile(sys["alarm-value"])
			Expect(err).To(BeNil())
			Expect(string(b)).To(Equal("0"))
		})
	})
})
