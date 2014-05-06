package gogadgets

import (
	"bitbucket.org/cswank/gogadgets/input"
	"bitbucket.org/cswank/gogadgets/models"
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

//Each input and output device has a config method that returns a models.Pin with
//the required fields poplulated with helpful values.
func GetTypes() map[string]models.ConfigHelper{
	t := input.Thermometer{}
	s := input.Switch{}
	g := output.GPIO{}
	h := output.Heater{}
	c := output.Cooler{}
	r := output.Recorder{}
	return map[string]models.ConfigHelper{
		"thermometer": t.Config(),
		"switch":      s.Config(),
		"gpio":     g.Config(),
		"heater":   h.Config(),
		"cooler":   c.Config(),
		"recorder": r.Config(),
	}
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
