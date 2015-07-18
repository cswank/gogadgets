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

var _ = Describe("Companies", func() {
	var ()
	BeforeEach(func() {

	})
	AfterEach(func() {
	})
	Describe("app", func() {
		It("starts up a gogadgets app", func() {
			port := 1024 + rand.Intn(65535-1024)
			p := &gogadgets.Gadget{
				Location:   "tank",
				Name:       "pump",
				OnCommand:  fmt.Sprintf("turn on %s %s", "tank", "pump"),
				OffCommand: fmt.Sprintf("turn off %s %s", "tank", "pump"),
				Output:     &FakeOutput{},
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
			}
			a := gogadgets.NewApp(cfg)
			a.AddGadget(p)
			a.AddGadget(s)
			go a.Start()
			time.Sleep(1 * time.Second)
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
