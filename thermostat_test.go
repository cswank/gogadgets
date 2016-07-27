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

type thermCase struct {
	temperature float64
	output      string
}

var _ = Describe("thermostat", func() {
	var (
		tmp   string
		sys   map[string]string
		therm gogadgets.OutputDevice
		pin   *gogadgets.Pin
		val   *gogadgets.Value
	)
	BeforeEach(func() {
		var err error
		tmp, err = ioutil.TempDir("", "")
		Expect(err).To(BeNil())
		sys = setupGPIOs(
			tmp,
			map[string]string{
				"heat": gogadgets.Pins["gpio"]["8"]["11"],
				"cool": gogadgets.Pins["gpio"]["8"]["12"],
			})
		gogadgets.GPIO_DEV_PATH = tmp
		gogadgets.GPIO_DEV_MODE = 0777

		pin = &gogadgets.Pin{
			Pins: map[string]gogadgets.Pin{
				"heat": {
					Type:      "gpio",
					Port:      "8",
					Pin:       "11",
					Direction: "out",
				},
				"cool": {
					Type:      "gpio",
					Port:      "8",
					Pin:       "12",
					Direction: "out",
				},
			},
			Args: map[string]interface{}{
				"sensor":  "my thermometer",
				"timeout": "0s",
			},
		}
		therm, err = gogadgets.NewThermostat(pin)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		os.RemoveAll(tmp)
	})
	Describe("heater", func() {

		BeforeEach(func() {
			val = &gogadgets.Value{
				Cmd:   "heat house to 70 F",
				Value: float64(70.0),
			}
		})

		It("sets up the gpio stuff correctly", func() {
			b, err := ioutil.ReadFile(sys["heat-direction"])
			Expect(err).To(BeNil())
			Expect(string(b)).To(Equal("out"))

			b, err = ioutil.ReadFile(sys["heat-value"])
			Expect(err).To(BeNil())
			Expect(string(b)).To(Equal("0"))
		})
		It("turns on", func() {
			Expect(therm.On(val)).To(BeNil())
			msg := &gogadgets.Message{
				Sender: "my thermometer",
				Value: gogadgets.Value{
					Value: 69.0,
				},
			}
			therm.Update(msg)
			b, err := ioutil.ReadFile(sys["heat-value"])
			Expect(err).To(BeNil())
			Expect(string(b)).To(Equal("1"))
		})
		It("turns off when the temperature is above range and back on when the temperature is below range", func() {
			Expect(therm.On(val)).To(BeNil())

			msg := &gogadgets.Message{
				Sender: "my thermometer",
				Value: gogadgets.Value{
					Value: 69.0,
				},
			}
			therm.Update(msg)

			b, err := ioutil.ReadFile(sys["heat-value"])
			Expect(err).To(BeNil())
			Expect(string(b)).To(Equal("1"))

			cases := []thermCase{
				{69.0, "1"},
				{70.0, "0"},
				{69.0, "1"},
				{29.0, "1"},
				{31.0, "1"},
				{69.0, "1"},
				{29.0, "1"},
				{70.0, "0"},
			}

			for _, c := range cases {
				msg := &gogadgets.Message{
					Sender: "my thermometer",
					Value: gogadgets.Value{
						Value: c.temperature,
					},
				}
				therm.Update(msg)
				b, err := ioutil.ReadFile(sys["heat-value"])
				Expect(err).To(BeNil())
				Expect(string(b)).To(Equal(c.output))
			}
		})
	})

	Describe("cooler", func() {

		var (
			val *gogadgets.Value
		)
		BeforeEach(func() {
			val = &gogadgets.Value{
				Cmd:   "cool house",
				Value: float64(70.0),
			}
		})
		It("turns on", func() {
			Expect(therm.On(val)).To(BeNil())
			msg := &gogadgets.Message{
				Sender: "my thermometer",
				Value: gogadgets.Value{
					Value: 70.1,
				},
			}
			therm.Update(msg)
			b, err := ioutil.ReadFile(sys["cool-value"])
			Expect(err).To(BeNil())
			Expect(string(b)).To(Equal("1"))
			Expect(therm.Status()).To(BeTrue())
		})

		It("gives the right status", func() {
			Expect(therm.Status()).To(BeFalse())
			Expect(therm.On(val)).To(BeNil())
			msg := &gogadgets.Message{
				Sender: "my thermometer",
				Value: gogadgets.Value{
					Value: 70.1,
				},
			}
			therm.Update(msg)
			Expect(therm.Status()).To(BeTrue())
		})

		It("turns off when the temperature is above range and back on when the temperature is below range", func() {
			Expect(therm.On(val)).To(BeNil())
			msg := &gogadgets.Message{
				Sender: "my thermometer",
				Value: gogadgets.Value{
					Value: 70.1,
				},
			}
			therm.Update(msg)

			b, err := ioutil.ReadFile(sys["cool-value"])
			Expect(err).To(BeNil())
			Expect(string(b)).To(Equal("1"))

			cases := []thermCase{
				{69.0, "0"},
				{70.0, "1"},
				{69.0, "0"},
				{29.0, "0"},
				{74.0, "1"},
				{31.0, "0"},
				{69.0, "0"},
				{29.0, "0"},
				{70.0, "1"},
			}

			for _, c := range cases {
				msg := &gogadgets.Message{
					Sender: "my thermometer",
					Value: gogadgets.Value{
						Value: c.temperature,
					},
				}
				therm.Update(msg)
				b, err := ioutil.ReadFile(sys["cool-value"])
				Expect(err).To(BeNil())
				Expect(string(b)).To(Equal(c.output))
			}
		})
	})

	Describe("temperature for other sensors", func() {

		It("does not react to other sensors", func() {
			val = &gogadgets.Value{
				Cmd:   "cool house",
				Value: float64(70.0),
			}
			Expect(therm.On(val)).To(BeNil())
			msg := &gogadgets.Message{
				Sender: "my thermometer",
				Value: gogadgets.Value{
					Value: 70.1,
				},
			}
			therm.Update(msg)
			b, err := ioutil.ReadFile(sys["cool-value"])
			Expect(err).To(BeNil())
			Expect(string(b)).To(Equal("1"))

			cases := []thermCase{
				{69.0, "1"},
				{150.0, "1"},
				{69.0, "1"},
				{129.0, "1"},
				{131.0, "1"},
				{69.0, "1"},
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
				therm.Update(msg)
				b, err := ioutil.ReadFile(sys["cool-value"])
				Expect(err).To(BeNil())
				Expect(string(b)).To(Equal(c.output))
			}
		})
	})
})

func setupGPIOs(pth string, pins map[string]string) map[string]string {
	sys := map[string]string{
		"export": path.Join(pth, "export"),
	}

	for k, pin := range pins {
		d := path.Join(pth, fmt.Sprintf("gpio%s", pin))
		Expect(os.Mkdir(d, 0777)).To(BeNil())
		sys[fmt.Sprintf("%s-direction", k)] = path.Join(d, "direction")
		sys[fmt.Sprintf("%s-edge", k)] = path.Join(d, "edge")
		sys[fmt.Sprintf("%s-value", k)] = path.Join(d, "value")
		sys[fmt.Sprintf("%s-active_low", k)] = path.Join(d, "active_low")
	}

	for _, v := range sys {
		Expect(ioutil.WriteFile(v, []byte(""), 0777)).To(BeNil())
	}
	return sys
}
