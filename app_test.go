package gogadgets_test

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"

	"github.com/cswank/gogadgets"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type fakeLogger struct {
	f bool
}

func (f *fakeLogger) Println(v ...interface{}) {}
func (f *fakeLogger) Fatal(v ...interface{})   { f.f = true }

func init() {
	rand.Seed(time.Now().Unix())
}

var _ = Describe("Companies", func() {
	var (
		port int
		lg   *fakeLogger
	)
	BeforeEach(func() {
		port = 1024 + rand.Intn(65535-1024)
		lg = &fakeLogger{}
	})
	AfterEach(func() {
	})
	Describe("app", func() {
		It("starts up a gogadgets app", func() {
			fo := &FakeOutput{}
			p := &gogadgets.Gadget{
				Location:   "tank",
				Name:       "pump",
				OnCommand:  fmt.Sprintf("turn on %s %s", "tank", "pump"),
				OffCommand: fmt.Sprintf("turn off %s %s", "tank", "pump"),
				Output:     fo,
				UID:        fmt.Sprintf("%s %s", "tank", "pump"),
			}
			location := "tank"
			name := "switch"
			poller := &FakePoller{}
			s := &gogadgets.Gadget{
				Location: location,
				Name:     name,
				Input: &gogadgets.Switch{
					GPIO:      poller,
					Value:     5.0,
					TrueValue: 5.0,
					Units:     "liters",
				},
				UID: fmt.Sprintf("%s %s", location, name),
			}
			cfg := &gogadgets.Config{
				Master:  true,
				Host:    "localhost",
				SubPort: port,
				PubPort: port + 1,
				Logger:  lg,
			}
			a := gogadgets.NewApp(cfg)
			a.AddGadget(p)
			a.AddGadget(s)

			input := make(chan gogadgets.Message)
			go a.GoStart(input)

			msg := gogadgets.Message{
				Type: "command",
				Body: "turn on tank pump",
			}

			Expect(fo.on).To(BeFalse())

			input <- msg

			Eventually(func() bool {
				return fo.on
			}).Should(BeTrue())
		})
		It("starts up swarm of gogadgets apps", func() {
			fo1 := &FakeOutput{}
			fo2 := &FakeOutput{}
			light1 := &gogadgets.Gadget{
				Location:   "living room",
				Name:       "light",
				OnCommand:  fmt.Sprintf("turn on %s %s", "living room", "light"),
				OffCommand: fmt.Sprintf("turn off %s %s", "living room", "light"),
				Output:     fo1,
				UID:        fmt.Sprintf("%s %s", "living room", "light"),
			}
			light2 := &gogadgets.Gadget{
				Location:   "kitchen",
				Name:       "light",
				OnCommand:  fmt.Sprintf("turn on %s %s", "kitchen", "light"),
				OffCommand: fmt.Sprintf("turn off %s %s", "kitchen", "light"),
				Output:     fo2,
				UID:        fmt.Sprintf("%s %s", "kitchen room", "light"),
			}

			cfg := &gogadgets.Config{
				Master:  true,
				Host:    "localhost",
				SubPort: port,
				PubPort: port + 1,
				Logger:  lg,
			}

			cfg2 := &gogadgets.Config{
				Master:  false,
				Host:    "localhost",
				SubPort: port,
				PubPort: port + 1,
				Logger:  lg,
			}

			a := gogadgets.NewApp(cfg)
			a.AddGadget(light1)
			a2 := gogadgets.NewApp(cfg2)
			a2.AddGadget(light2)

			input := make(chan gogadgets.Message)
			go a.GoStart(input)
			time.Sleep(100 * time.Millisecond)
			go a2.Start()

			Expect(fo1.on).To(BeFalse())
			Expect(fo2.on).To(BeFalse())

			msg := gogadgets.Message{
				Sender: "the test",
				Type:   "command",
				Body:   "turn on living room light",
			}

			time.Sleep(500 * time.Millisecond)

			input <- msg

			Eventually(func() bool {
				return fo1.on
			}).Should(BeTrue())
			Expect(fo2.on).To(BeFalse())

			msg = gogadgets.Message{
				Sender: "the test",
				Type:   "command",
				Body:   "turn on kitchen light",
			}

			Expect(fo1.on).To(BeTrue())
			Expect(fo2.on).To(BeFalse())

			input <- msg

			Eventually(func() bool {
				return fo2.on
			}).Should(BeTrue())
		})
		It("loads a json config file", func() {
			s := `{
    "gadgets": [
        {
            "location": "front yard",
            "name": "sprinklers",
            "pin": {
                "type": "gpio",
                "port": "8",
                "pin": "10",
                "direction": "out"
            }
        },
        {
            "location": "front garden",
            "name": "sprinklers",
            "pin": {
                "type": "gpio",
                "port": "8",
                "pin": "11",
                "direction": "out"
            }
        },
        {
            "location": "sidewalk",
            "name": "sprinklers",
            "pin": {
                "type": "gpio",
                "port": "8",
                "pin": "12",
                "direction": "out"
            }
        },
        {
            "location": "back yard",
            "name": "sprinklers",
            "pin": {
                "type": "gpio",
                "port": "8",
                "pin": "14",
                "direction": "out"
            }
        },
        {
            "location": "back garden",
            "name": "sprinklers",
            "pin": {
                "type": "gpio",
                "port": "8",
                "pin": "15",
                "direction": "out"
            }
        }
    ]
}
`
			f, _ := ioutil.TempFile("", "")
			f.Write([]byte(s))
			f.Close()
			cfg := gogadgets.GetConfig(f.Name())
			os.Remove(f.Name())
			Expect(len(cfg.Gadgets)).To(Equal(5))
		})
	})
})
