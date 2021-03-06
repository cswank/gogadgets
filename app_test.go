package gogadgets_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/cswank/gogadgets"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func init() {
	rand.Seed(time.Now().Unix())
}

var _ = Describe("gogadgets", func() {
	var (
		port int
	)
	BeforeEach(func() {
		port = 1024 + rand.Intn(65535-1024)
	})

	Describe("app", func() {
		It("starts up a gogadgets app", func() {
			fo := &FakeOutput{}
			p := &gogadgets.Gadget{
				Location:    "tank",
				Name:        "pump",
				OnCommands:  []string{fmt.Sprintf("turn on %s %s", "tank", "pump")},
				OffCommands: []string{fmt.Sprintf("turn off %s %s", "tank", "pump")},
				Output:      fo,
				UID:         fmt.Sprintf("%s %s", "tank", "pump"),
			}
			location := "tank"
			name := "switch"
			poller := &FakePoller{}
			s := &gogadgets.Gadget{
				Location: location,
				Name:     name,
				Input: &gogadgets.Switch{
					GPIO: poller,
				},
				UID: fmt.Sprintf("%s %s", location, name),
			}
			cfg := &gogadgets.Config{
				Master: "",
				Host:   "localhost",
				Port:   port,
			}
			a := gogadgets.New(cfg, p, s)

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
				Location:    "living room",
				Name:        "light",
				OnCommands:  []string{fmt.Sprintf("turn on %s %s", "living room", "light")},
				OffCommands: []string{fmt.Sprintf("turn off %s %s", "living room", "light")},
				Output:      fo1,
				UID:         fmt.Sprintf("%s %s", "living room", "light"),
			}
			light2 := &gogadgets.Gadget{
				Location:    "kitchen",
				Name:        "light",
				OnCommands:  []string{fmt.Sprintf("turn on %s %s", "kitchen", "light")},
				OffCommands: []string{fmt.Sprintf("turn off %s %s", "kitchen", "light")},
				Output:      fo2,
				UID:         fmt.Sprintf("%s %s", "kitchen room", "light"),
			}

			cfg := &gogadgets.Config{
				Master: "",
				Host:   "",
				Port:   port,
			}

			cfg2 := &gogadgets.Config{
				Master: fmt.Sprintf("http://localhost:%d", port),
				Host:   fmt.Sprintf("http://localhost:%d", port+1),
				Port:   port + 1,
			}

			a := gogadgets.New(cfg, light1)
			a2 := gogadgets.New(cfg2, light2)

			input := make(chan gogadgets.Message)
			go a.GoStart(input)
			time.Sleep(100 * time.Millisecond)
			go a2.Start()

			Eventually(func() bool {
				r, err := http.Get(fmt.Sprintf("http://localhost:%d/clients", port))
				if err != nil || r.StatusCode != http.StatusOK {
					return false
				}
				var c map[string]string
				dec := json.NewDecoder(r.Body)
				dec.Decode(&c)
				r.Body.Close()
				return len(c) > 0
			}).Should(BeTrue())

			Expect(fo1.on).To(BeFalse())
			Expect(fo2.on).To(BeFalse())

			msg := gogadgets.Message{
				Sender: "the test",
				Type:   "command",
				Body:   "turn on living room light",
			}

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
