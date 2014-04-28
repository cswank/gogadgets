package gogadgets

import (
	"bitbucket.org/cswank/gogadgets/input"
	"bitbucket.org/cswank/gogadgets/output"
)

type AppFactory struct {
	inputFactories  map[string]input.InputDeviceFactory
	outputFactories map[string]output.OutputDeviceFactory
}

func NewAppFactory() *AppFactory {
	a := &AppFactory{
		inputFactories: map[string]input.InputDeviceFactory{
			"thermometer": input.NewThermometer,
			"switch":      input.NewSwitch,
		},
		outputFactories: map[string]output.OutputDeviceFactory{
			"gpio":     output.NewGPIO,
			"heater":   output.NewHeater,
			"cooler":   output.NewCooler,
			"recorder": output.NewRecorder,
		},
	}
	return a
}

func (f *AppFactory) RegisterInputFactory(name string, factory input.InputDeviceFactory) {
	f.inputFactories[name] = factory
}

func (f *AppFactory) RegisterOutputFactory(name string, factory output.OutputDeviceFactory) {
	f.outputFactories[name] = factory
}

func (f *AppFactory) GetApp() (a *App, err error) {
	return a, err
}
