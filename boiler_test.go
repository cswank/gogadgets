package gogadgets_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/cswank/gogadgets"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("boiler", func() {
	var (
		tmp    string
		sys    map[string]string
		boiler gogadgets.OutputDevice
		pin    *gogadgets.Pin
	)
	BeforeEach(func() {
		var err error
		tmp, err = ioutil.TempDir("", "")
		Expect(err).To(BeNil())
		sys = setupGPIO(tmp, gogadgets.Pins["gpio"]["8"]["11"])
		gogadgets.GPIO_DEV_PATH = tmp
		gogadgets.GPIO_DEV_MODE = 0777
	})

	AfterEach(func() {
		os.RemoveAll(tmp)
	})
	Describe("heater", func() {
		BeforeEach(func() {
			pin = &gogadgets.Pin{
				Port:      "8",
				Pin:       "11",
				Direction: "out",
				Args: map[string]interface{}{
					"type":   "heater",
					"high":   150.0,
					"low":    130.0,
					"sensor": "my thermometer",
				},
			}
			var err error
			boiler, err = gogadgets.NewBoiler(pin)
			Expect(err).To(BeNil())
		})

		It("sets up the gpio stuff correctly", func() {
			b, err := ioutil.ReadFile(sys["direction"])
			Expect(err).To(BeNil())
			Expect(string(b)).To(Equal("out"))

			b, err = ioutil.ReadFile(sys["value"])
			Expect(err).To(BeNil())
			Expect(string(b)).To(Equal("0"))
		})
		It("turns on", func() {
			Expect(boiler.On(nil)).To(BeNil())
			b, err := ioutil.ReadFile(sys["value"])
			Expect(err).To(BeNil())
			Expect(string(b)).To(Equal("1"))
		})
		It("turns off when the temperature is above range and back on when the temperature is below range", func() {
			Expect(boiler.On(nil)).To(BeNil())
			b, err := ioutil.ReadFile(sys["value"])
			Expect(err).To(BeNil())
			Expect(string(b)).To(Equal("1"))

			cases := []thermCase{
				{149.0, "1"},
				{150.0, "0"},
				{149.0, "0"},
				{129.0, "1"},
				{131.0, "1"},
				{149.0, "1"},
				{129.0, "1"},
				{150.0, "0"},
			}

			for _, c := range cases {
				msg := &gogadgets.Message{
					Sender: "my thermometer",
					Value: gogadgets.Value{
						Value: c.temperature,
					},
				}
				boiler.Update(msg)
				b, err := ioutil.ReadFile(sys["value"])
				Expect(err).To(BeNil())
				Expect(string(b)).To(Equal(c.output))
			}
		})
	})

	Describe("cooler", func() {
		BeforeEach(func() {
			pin = &gogadgets.Pin{
				Port:      "8",
				Pin:       "11",
				Direction: "out",
				Args: map[string]interface{}{
					"type":   "cooler",
					"high":   150.0,
					"low":    130.0,
					"sensor": "my thermometer",
				},
			}
			var err error
			boiler, err = gogadgets.NewBoiler(pin)
			Expect(err).To(BeNil())
		})
		It("turns on", func() {
			Expect(boiler.On(nil)).To(BeNil())
			b, err := ioutil.ReadFile(sys["value"])
			Expect(err).To(BeNil())
			Expect(string(b)).To(Equal("1"))
		})
		It("turns off when the temperature is above range and back on when the temperature is below range", func() {
			Expect(boiler.On(nil)).To(BeNil())
			b, err := ioutil.ReadFile(sys["value"])
			Expect(err).To(BeNil())
			Expect(string(b)).To(Equal("1"))

			cases := []thermCase{
				{149.0, "1"},
				{150.0, "1"},
				{149.0, "1"},
				{129.0, "0"},
				{131.0, "0"},
				{149.0, "0"},
				{129.0, "0"},
				{150.0, "1"},
			}

			for _, c := range cases {
				msg := &gogadgets.Message{
					Sender: "my thermometer",
					Value: gogadgets.Value{
						Value: c.temperature,
					},
				}
				boiler.Update(msg)
				b, err := ioutil.ReadFile(sys["value"])
				Expect(err).To(BeNil())
				Expect(string(b)).To(Equal(c.output))
			}
		})
	})

	Describe("temperature for other sensors", func() {
		BeforeEach(func() {
			pin = &gogadgets.Pin{
				Port:      "8",
				Pin:       "11",
				Direction: "out",
				Args: map[string]interface{}{
					"type":   "cooler",
					"high":   150.0,
					"low":    130.0,
					"sensor": "my thermometer",
				},
			}
			var err error
			boiler, err = gogadgets.NewBoiler(pin)
			Expect(err).To(BeNil())
		})

		It("does not react to other sensors", func() {
			Expect(boiler.On(nil)).To(BeNil())
			b, err := ioutil.ReadFile(sys["value"])
			Expect(err).To(BeNil())
			Expect(string(b)).To(Equal("1"))

			cases := []thermCase{
				{149.0, "1"},
				{150.0, "1"},
				{149.0, "1"},
				{129.0, "1"},
				{131.0, "1"},
				{149.0, "1"},
				{129.0, "1"},
				{150.0, "1"},
			}

			for _, c := range cases {
				msg := &gogadgets.Message{
					Sender: "that other thermometer",
					Value: gogadgets.Value{
						Value: c.temperature,
					},
				}
				boiler.Update(msg)
				b, err := ioutil.ReadFile(sys["value"])
				Expect(err).To(BeNil())
				Expect(string(b)).To(Equal(c.output))
			}
		})
	})
})

func setupGPIO(pth string, pin string) map[string]string {
	d := path.Join(pth, fmt.Sprintf("gpio%s", pin))
	sys := map[string]string{
		"export":     path.Join(pth, "export"),
		"direction":  path.Join(d, "direction"),
		"edge":       path.Join(d, "edge"),
		"value":      path.Join(d, "value"),
		"active_low": path.Join(d, "active_low"),
	}
	Expect(os.Mkdir(d, 0777)).To(BeNil())
	for _, v := range sys {
		Expect(ioutil.WriteFile(v, []byte(""), 0777)).To(BeNil())
	}
	return sys
}
