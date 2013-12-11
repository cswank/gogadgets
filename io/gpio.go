package io

import (
	"bitbucket.com/cswank/gogadgets.pins"
	"io"
	"io/ioutil"
)

type GPOutput struct {
	OutputDevice
	Pin *pins.GPIO
	exportPath string
}

func NewGPOutput(port, pin string) (*GPOutput, error) {
	port, ok := pins.GPIO[port]
	if !ok {
		return nil, errors.New(fmt.Sprintf("no such port: %s", port))
	}
	export, ok := port[pin]
	if !ok {
		return nil, errors.New(fmt.Sprintf("no such pin: %s", pin))
	}
	g := &GPOutput{
		Pin: pin,
		exportPath: "/sys/class/gpio/export",
		directionPath: fmt.Sprintf("/sys/class/gpio/gpio%s/direction", export),
		valuePath: fmt.Sprintf("/sys/class/gpio/gpio%s/value", export),
	}
	err := g.Init()
	return g, err
}

func (g *GPOutput) Init() error {
	err := g.writeValue(g.exportPath, g.Pin.Export)
	if err == nil {
		err := g.writeValue(g.directionPath, "out")
		if err == nil {
			err := g.writeValue(g.valuePath, "0")
		}
	}
}

func (g *GPOutput) On() error {
	return g.writeValue(g.valuePath, "1")
}

func (g *GPOutput) Off() error {
	return g.writeValue(g.valuePath, "0")
}

func (g *GPOutput) griteValue(path, value string) error {
	return ioutil.WriteFile(path, []byte(value), io.DeviceFile)
}
