package gogadgets_test

import (
	"sort"
	"sync"

	"github.com/cswank/gogadgets"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	serial "go.bug.st/serial.v1"
)

type FakePort struct {
	serial.Port
	msg  []byte
	i    int
	lock sync.Mutex
	end  int
}

func (f *FakePort) SetMode(mode *serial.Mode) error {
	return nil
}

func (f *FakePort) BetModemStatusBits(mode *serial.Mode) error {
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
		io     string
		fp     *FakePort
	)

	BeforeEach(func() {
		fp = &FakePort{}
	})

	JustBeforeEach(func() {
		fp.setMsg(packet)

		pin := &gogadgets.Pin{
			Args: map[string]interface{}{
				"port":  "fake-port",
				"xbees": io,
			},
		}

		val = make(chan gogadgets.Value)
		msg = make(chan gogadgets.Message)

		var err error
		xbee, err = gogadgets.NewXBee(pin, gogadgets.XBeeSerialPort(fp))
		Expect(err).To(BeNil())

		go xbee.Start(msg, val)
	})

	Context("adc", func() {

		BeforeEach(func() {
			packet = []byte{0x7E, 0x00, 0x14, 0x92, 0x00, 0x13, 0xA2, 0x00, 0x40, 0x4C, 0x0E, 0xBE, 0x61, 0x59, 0x01, 0x01, 0x00, 0x00, 0x03, 0x02, 0x2F, 0x01, 0xFE, 0x71}
			io = `{"13a200404c0ebe": {"location": "garden", "adc": {"adc0": "tmp36","adc1": "tmp36"}}}`
		})

		It("reports the xbee analog values", func() {
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

	Context("dio", func() {

		BeforeEach(func() {
			packet = []byte{0x7E, 0x00, 0x12, 0x92, 0x00, 0x13, 0xA2, 0x00, 0x40, 0x4C, 0x0E, 0xBE, 0x61, 0x59, 0x01, 0x01, 0x00, 0x18, 0x00, 0x00, 0x10, 0x7c}
			io = `{"13a200404c0ebe": {"location": "garden", "dio": {"dio3": "gate", "dio4": "door"}}}`
		})

		It("reports the xbee digital io", func() {
			vals := make([]gogadgets.Value, 2)
			vals[0] = <-val
			vals[1] = <-val

			sort.Sort(byName(vals))
			Expect(vals[0].GetName()).To(Equal("door"))
			Expect(vals[0].Value.(bool)).To(BeTrue())
			Expect(vals[1].GetName()).To(Equal("gate"))
			Expect(vals[1].Value.(bool)).To(BeFalse())
		})
	})

	Context("both", func() {

		BeforeEach(func() {
			packet = []byte{0x7E, 0x00, 0x16, 0x92, 0x00, 0x13, 0xA2, 0x00, 0x40, 0x4C, 0x0E, 0xBE, 0x61, 0x59, 0x01, 0x01, 0x00, 0x18, 0x03, 0x00, 0x10, 0x02, 0x2F, 0x01, 0xFE, 0x49}
			io = `{"13a200404c0ebe": {"location": "garden", "dio": {"dio3": "gate", "dio4": "door"}, "adc": {"adc0": "tmp36","adc1": "moisture"}}}`
		})

		It("reports the xbee analog and digital values", func() {
			vals := make([]gogadgets.Value, 4)
			vals[0] = <-val
			vals[1] = <-val
			vals[2] = <-val
			vals[3] = <-val

			sort.Sort(byName(vals))
			Expect(vals[0].GetName()).To(Equal("door"))
			Expect(vals[0].Value.(bool)).To(BeTrue())
			Expect(vals[1].GetName()).To(Equal("gate"))
			Expect(vals[1].Value.(bool)).To(BeFalse())
			Expect(vals[2].GetName()).To(Equal("moisture"))
			Expect(vals[2].Value.(float64)).To(BeNumerically("~", 4.9824046920821115, 0.001))
			Expect(vals[3].GetName()).To(Equal("temperature"))
			Expect(vals[3].Value.(float64)).To(BeNumerically("~", 60.0293255, 0.001))
		})
	})
})

type byName []gogadgets.Value

func (a byName) Len() int           { return len(a) }
func (a byName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byName) Less(i, j int) bool { return a[i].GetName() < a[j].GetName() }
