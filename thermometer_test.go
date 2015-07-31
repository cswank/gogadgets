package gogadgets_test

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/cswank/gogadgets"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func init() {
	rand.Seed(time.Now().Unix())
}

var _ = Describe("Thermometer", func() {
	var (
		tmp  string
		p    string
		f    *os.File
		pin  *gogadgets.Pin
		lock sync.Mutex
	)
	BeforeEach(func() {
		lock = sync.Mutex{}
		var err error
		tmp, err = ioutil.TempDir("", "")
		Expect(err).To(BeNil())
		pin = &gogadgets.Pin{
			OneWirePath: fmt.Sprintf("%s/%%s", tmp),
			OneWireId:   "fakeid",
			Sleep:       10 * time.Millisecond,
			Lock:        lock,
		}
		p = fmt.Sprintf(pin.OneWirePath, pin.OneWireId)
		f, err = os.Create(p)
		Expect(err).To(BeNil())
		f.Write([]byte("3d 01 4b 46 7f ff 03 10 6d : crc=6d YES\n3d 01 4b 46 7f ff 03 10 6d t=19812\n"))
		f.Close()
	})
	AfterEach(func() {
		os.RemoveAll(tmp)
	})
	Describe("NewThermometer", func() {
		It("creates a thermometer", func() {
			t, err := gogadgets.NewThermometer(pin)
			Expect(err).To(BeNil())
			Expect(t).ToNot(BeNil())
		})
	})
	Describe("Start", func() {
		It("reads the temperature", func() {
			t, _ := gogadgets.NewThermometer(pin)
			in := make(chan gogadgets.Value)
			out := make(chan gogadgets.Message)
			go t.Start(out, in)
			val := <-in
			Expect(val.Value).To(Equal(19.812))
		})
		It("updates the temperature if the change is reasonable", func() {
			t, _ := gogadgets.NewThermometer(pin)
			in := make(chan gogadgets.Value)
			out := make(chan gogadgets.Message)
			go t.Start(out, in)
			val := <-in
			Expect(val.Value).To(Equal(19.812))
			lock.Lock()
			var err error
			f, err = os.Create(p)
			Expect(err).To(BeNil())
			f.Write([]byte("3d 01 4b 46 7f ff 03 10 6d : crc=6d YES\n3d 01 4b 46 7f ff 03 10 6d t=20008\n"))
			f.Close()
			lock.Unlock()
			Eventually(func() float64 {
				val = <-in
				return val.Value.(float64)
			}).Should(Equal(20.008))
		})
		It("keeps the old temperature if the change is too large", func() {
			t, _ := gogadgets.NewThermometer(pin)
			in := make(chan gogadgets.Value)
			out := make(chan gogadgets.Message)
			go t.Start(out, in)
			val := <-in
			lock.Lock()
			Expect(val.Value).To(Equal(19.812))
			var err error
			f, err = os.Create(p)
			Expect(err).To(BeNil())
			f.Write([]byte("3d 01 4b 46 7f ff 03 10 6d : crc=6d YES\n3d 01 4b 46 7f ff 03 10 6d t=30812\n"))
			f.Close()
			lock.Unlock()
			time.Sleep(20 * time.Millisecond)
			v := t.GetValue()
			Expect(v.Value).To(Equal(19.812))
		})
	})
})
