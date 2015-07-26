package gogadgets_test

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/cswank/gogadgets"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func init() {
	rand.Seed(time.Now().Unix())
}

var _ = Describe("gadgets", func() {
	var (
		port int
	)
	BeforeEach(func() {
		port = 1024 + rand.Intn(65535-1024)
	})
	AfterEach(func() {
	})
	Describe("commands", func() {
		It("parses a command", func() {
			val, unit, err := gogadgets.ParseCommand("fill tank to 5 liters")
			Expect(err).To(BeNil())
			Expect(val).To(Equal(5.0))
			Expect(unit).To(Equal("liters"))
		})
		It("parses an off command", func() {
			val, unit, err := gogadgets.ParseCommand("turn on lab led to -50 %")
			Expect(err).To(BeNil())
			Expect(val).To(Equal(-50.0))
			Expect(unit).To(Equal("%"))
		})
		It("parses a time command", func() {
			val, unit, err := gogadgets.ParseCommand("turn on lab led for 1.1 minutes")
			Expect(err).To(BeNil())
			Expect(val).To(Equal(1.1))
			Expect(unit).To(Equal("minutes"))
		})
	})
	Describe("start", func() {
		It("starts a gadget and turns it on and off", func() {
			location := "lab"
			name := "led"
			g := gogadgets.Gadget{
				Location:   location,
				Name:       name,
				Direction:  "output",
				OnCommand:  fmt.Sprintf("turn on %s %s", location, name),
				OffCommand: fmt.Sprintf("turn off %s %s", location, name),
				Output:     &FakeOutput{},
				UID:        fmt.Sprintf("%s %s", location, name),
			}
			input := make(chan gogadgets.Message)
			output := make(chan gogadgets.Message)
			go g.Start(input, output)
			update := <-output
			Expect(update.Value.Value).To(BeFalse())
			msg := gogadgets.Message{
				Type: "command",
				Body: "turn on lab led",
			}
			input <- msg
			update = <-output
			Expect(update.Value.Value).To(BeTrue())
			Expect(update.Value.Output).To(BeTrue())

			msg = gogadgets.Message{
				Type: "command",
				Body: "turn off lab led",
			}
			input <- msg
			update = <-output
			Expect(update.Value.Value).To(BeFalse())
			Expect(update.Value.Output).To(BeFalse())

			msg = gogadgets.Message{
				Type: "command",
				Body: "shutdown",
			}
			input <- msg
			update = <-output
		})
		It("starts a gadgets that has a trigger", func() {
			location := "tank"
			name := "valve"
			g := gogadgets.Gadget{
				Location:   location,
				Name:       name,
				Operator:   ">=",
				OnCommand:  fmt.Sprintf("fill %s", location),
				OffCommand: fmt.Sprintf("stop filling %s", location),
				Output:     &FakeOutput{},
				UID:        fmt.Sprintf("%s %s", location, name),
			}
			input := make(chan gogadgets.Message)
			output := make(chan gogadgets.Message)
			go g.Start(input, output)
			update := <-output
			Expect(update.Value.Value).To(BeFalse())

			msg := gogadgets.Message{
				Type: "command",
				Body: "fill tank to 4.4 liters",
			}
			input <- msg
			update = <-output
			Expect(update.Value.Value).To(BeTrue())

			//make a message that should trigger the trigger and stop the device
			msg = gogadgets.Message{
				Sender:   "tank volume",
				Type:     gogadgets.UPDATE,
				Location: "tank",
				Name:     "volume",
				Value: gogadgets.Value{
					Units: "liters",
					Value: 4.4,
				},
			}
			input <- msg
			update = <-output
			Expect(update.Value.Value).To(BeFalse())

		})
		It("starts a gadget that has a time trigger", func() {
			location := "lab"
			name := "led"
			g := gogadgets.Gadget{
				Location:   location,
				Name:       name,
				OnCommand:  "turn on lab led",
				Operator:   ">=",
				OffCommand: "turn off lab led",
				Output:     &FakeOutput{},
				UID:        fmt.Sprintf("%s %s", location, name),
			}
			input := make(chan gogadgets.Message)
			output := make(chan gogadgets.Message)
			go g.Start(input, output)
			update := <-output
			Expect(update.Value.Value).To(BeFalse())

			msg := gogadgets.Message{
				Type: "command",
				Body: "turn on lab led for 0.01 seconds",
			}
			input <- msg
			update = <-output
			Expect(update.Value.Value).To(BeTrue())

			//wait for a second
			update = <-output
			Expect(update.Value.Value).To(BeFalse())
		})

		It("starts a gadget with a time trigger that gets interrupted", func() {
			location := "lab"
			name := "led"
			g := gogadgets.Gadget{
				Location:   location,
				Name:       name,
				OnCommand:  "turn on lab led",
				OffCommand: "turn off lab led",
				Output:     &FakeOutput{},
				UID:        fmt.Sprintf("%s %s", location, name),
			}
			input := make(chan gogadgets.Message)
			output := make(chan gogadgets.Message)
			go g.Start(input, output)
			update := <-output
			Expect(update.Value.Value).To(BeFalse())
			msg := gogadgets.Message{
				Type: "command",
				Body: "turn on lab led for 30 seconds",
			}
			input <- msg
			update = <-output
			Expect(update.Value.Value).To(BeTrue())

			msg = gogadgets.Message{
				Type: "command",
				Body: "turn on lab led",
			}
			input <- msg

			msg = gogadgets.Message{
				Type: "update",
				Body: "",
			}
			input <- msg

			msg = gogadgets.Message{
				Type: "command",
				Body: "turn off lab led",
			}
			input <- msg
			update = <-output
			Expect(update.Value.Value).To(BeFalse())
		})
		It("starts a switch", func() {
			location := "lab"
			name := "switch"
			poller := &FakePoller{}
			s := &gogadgets.Switch{
				GPIO:      poller,
				Value:     5.0,
				TrueValue: 5.0,
				Units:     "liters",
			}
			g := gogadgets.Gadget{
				Location: location,
				Name:     name,
				Input:    s,
				UID:      fmt.Sprintf("%s %s", location, name),
			}
			input := make(chan gogadgets.Message)
			output := make(chan gogadgets.Message)
			go g.Start(input, output)
			val := <-output
			Expect(val.Value.Value.(float64)).To(Equal(5.0))
			val = <-output
			Expect(val.Value.Value.(float64)).To(Equal(0.0))
		})
	})
})
