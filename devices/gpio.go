package devices

import (
	"fmt"
	"bitbucket.com/cswank/gogadgets/utils"
	"os"
	"errors"
	"io/ioutil"
)

type GPIO struct {
	OutputDevice
	InputDevice
	units string
	export string
	exportPath string
	directionPath string
	valuePath string
	direction string
	edge string
}

func NewGPIO(pin *Pin) (*GPIO, error) {
	portMap, ok := Pins["gpio"][pin.Port]
	if !ok {
		return nil, errors.New(fmt.Sprintf("no such port: %s", pin.Port))
	}
	export, ok := portMap[pin.Pin]
	if !ok {
		return nil, errors.New(fmt.Sprintf("no such pin: %s", pin.Pin))
	}
	g := &GPIO{
		export: export,
		exportPath: "/sys/class/gpio/export",
		directionPath: fmt.Sprintf("/sys/class/gpio/gpio%s/direction", export),
		valuePath: fmt.Sprintf("/sys/class/gpio/gpio%s/value", export),
		direction: pin.Direction,
		edge: pin.Edge,
	}
	err := g.Init()
	return g, err
}

func (g *GPIO) Init() error {
	var err error
	if !utils.FileExists(g.directionPath) {
		err = g.writeValue(g.exportPath, g.export)
	}
	if err == nil {
		err = g.writeValue(g.directionPath, "out")
		if err == nil {
			err = g.writeValue(g.valuePath, "0")
		}
	}
	return err
}

func (g *GPIO) On() error {
	return g.writeValue(g.valuePath, "1")
}

func (g *GPIO) Status() bool {
	data, err := ioutil.ReadFile(g.valuePath)
	return err == nil && string(data) == "1\n"
}

func (g *GPIO) Off() error {
	return g.writeValue(g.valuePath, "0")
}

func (g *GPIO) writeValue(path, value string) error {
	return ioutil.WriteFile(path, []byte(value), os.ModeDevice)
}
