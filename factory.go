package gogadgets

import (

)

type AppFactory struct {
	inputFactories map[string] InputDeviceFactory
	outputFactories map[string] OutputDeviceFactory
}

func NewAppFactory() *AppFactory {
	a := &AppFactory{
		inputFactories: map[string] InputDeviceFactory{
			"thermometer": NewThermometer,
			"switch": NewSwitch,
		},
		outputFactories: map[string] OutputDeviceFactory{
			"gpio": NewGPIO,
			"heater": NewHeater,
			"cooler": NewCooler,
		},
	}
	return a
}

func (f *AppFactory) RegisterInputFactory(name string, factory InputDeviceFactory) {
	f.inputFactories[name] = factory
}

func (f *AppFactory) RegisterOutputFactory(name string, factory OutputDeviceFactory) {
	f.outputFactories[name] = factory
}

func (f *AppFactory) GetApp() (a *App, err error) {
	return a, err
}
