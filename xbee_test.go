package gogadgets_test

import (
	"sort"
	"sync"

	"github.com/cswank/gogadgets"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	serial "go.bug.st/serial.v1"
)

var (
	fp *FakePort
)

type FakePort struct {
	msg  []byte
	i    int
	lock sync.Mutex
	end  int
}

func NewFakePort(p string, mode *serial.Mode) (serial.Port, error) {
	return fp, nil
}

func (f *FakePort) SetMode(mode *serial.Mode) error {
	return nil
}

func (f *FakePort) setMsg(msg []byte) {
	f.msg = msg
	f.end = len(msg)
}

func (f *FakePort) Read(p []byte) (n int, err error) {
	f.lock.Lock()
	copy(p, f.msg[f.i:])
	f.i += len(p)
	f.lock.Unlock()
	if f.i == f.end {
		f.lock.Lock()
	}
	return len(p), nil
}

func (f *FakePort) Write(p []byte) (n int, err error) {
	return 0, nil
}

func (f *FakePort) Close() error {
	return nil
}

var _ = Describe("xbee", func() {
	var (
		xbee   gogadgets.InputDevice
		packet []byte
		msg    chan gogadgets.Message
		val    chan gogadgets.Value
	)

	BeforeEach(func() {
		fp = &FakePort{}
		gogadgets.Init(NewFakePort)

		pin := &gogadgets.Pin{
			Args: map[string]interface{}{
				"port":  "fake-port",
				"xbees": `{"13a200404c0ebe": {"location": "garden", "pins": {"adc0": "tmp36","adc1": "tmp36"}}}`,
			},
		}

		val = make(chan gogadgets.Value)
		msg = make(chan gogadgets.Message)

		var err error
		xbee, err = gogadgets.NewXBee(pin)
		Expect(err).To(BeNil())

		go xbee.Start(msg, val)

	})

	JustBeforeEach(func() {
		fp.setMsg(packet)
	})

	AfterEach(func() {

	})

	Describe("acd", func() {

		BeforeEach(func() {
			packet = []byte{0x7E, 0x00, 0x16, 0x92, 0x00, 0x13, 0xA2, 0x00, 0x40, 0x4C, 0x0E, 0xBE, 0x61, 0x59, 0x01, 0x01, 0x00, 0x18, 0x03, 0x00, 0x10, 0x02, 0x2F, 0x01, 0xFE, 0x49}
		})

		It("reports the xbee update", func() {
			vals := make([]float64, 2)
			v := <-val
			vals[0] = v.Value.(float64)
			v = <-val
			vals[1] = v.Value.(float64)
			sort.Float64s(vals)
			Expect(vals[0]).To(BeNumerically("~", 49.68328445747801, 0.001))
			Expect(vals[1]).To(BeNumerically("~", 60.0293255, 0.001))
		})
	})
})
