package gogadgets_test

import (
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
				"xbees": `{"13a200409825c1": {"location": "garden", "pins": {"adc0": "moisture","adc1": "tmp36"}}}`,
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
			packet = []byte{126, 0, 18, 146, 0, 19, 162, 0, 64, 152, 37, 193, 222, 186, 1, 1, 0, 0, 2, 2, 13, 79}
		})

		It("reports the xbee update", func() {
			v := <-val
			Expect(v.Value.(float64)).To(BeNumerically("~", 615.8357, 0.001))
		})
	})
})
