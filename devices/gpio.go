package devices

import (
	"fmt"
	"bitbucket.com/cswank/gogadgets/pins"
	"bitbucket.com/cswank/gogadgets/utils"
	"os"
	"errors"
	"io/ioutil"
)

type GPOutput struct {
	OutputDevice
	export string
	exportPath string
	directionPath string
	valuePath string
}

func NewGPOutput(pin *pins.Pin) (*GPOutput, error) {
	portMap, ok := pins.GPIO[pin.Port]
	if !ok {
		return nil, errors.New(fmt.Sprintf("no such port: %s", pin.Port))
	}
	export, ok := portMap[pin.Pin]
	if !ok {
		return nil, errors.New(fmt.Sprintf("no such pin: %s", pin.Pin))
	}
	g := &GPOutput{
		export: export,
		exportPath: "/sys/class/gpio/export",
		directionPath: fmt.Sprintf("/sys/class/gpio/gpio%s/direction", export),
		valuePath: fmt.Sprintf("/sys/class/gpio/gpio%s/value", export),
	}
	err := g.Init()
	return g, err
}

func (g *GPOutput) Init() error {
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

func (g *GPOutput) On() error {
	return g.writeValue(g.valuePath, "1")
}

func (g *GPOutput) Off() error {
	return g.writeValue(g.valuePath, "0")
}

func (g *GPOutput) writeValue(path, value string) error {
	return ioutil.WriteFile(path, []byte(value), os.ModeDevice)
}
