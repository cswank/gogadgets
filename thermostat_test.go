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
	)
	BeforeEach(func() {
		var err error
		tmp, err = ioutil.TempDir("", "")
		Expect(err).To(BeNil())
		sys = setupGPIO(tmp, gogadgets.Pins["gpio"]["8"]["11"])
		gogadgets.GPIO_DEV_PATH = tmp
		gogadgets.GPIO_DEV_MODE = 0777
		pin := &gogadgets.Pin{
			Port:      "8",
			Pin:       "11",
			Direction: "out",
			Args: map[string]interface{}{
				"high": 150.0,
				"low":  130.0,
			},
		}
		therm, err = gogadgets.NewThermostat(pin)
		Expect(err).To(BeNil())
	})

	AfterEach(func() {
		os.RemoveAll(tmp)
	})
	Describe("heater", func() {
		FIt("turns on", func() {
			Expect(therm.On(nil)).To(BeNil())
			b, err := ioutil.ReadFile(sys["value"])
			Expect(err).To(BeNil())
			Expect(string(b)).To(Equal("1"))
		})
		FIt("turns off when the temperature is above range and back on when the temperature is below range", func() {
			Expect(therm.On(nil)).To(BeNil())
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
					Value: gogadgets.Value{
						Value: c.temperature,
					},
				}
				therm.Update(msg)
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
