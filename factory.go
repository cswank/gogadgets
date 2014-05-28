package gogadgets

type AppFactory struct {
	inputFactories  map[string]InputDeviceFactory
	outputFactories map[string]OutputDeviceFactory
}

func NewAppFactory() *AppFactory {
	a := &AppFactory{
		inputFactories: map[string]InputDeviceFactory{
			"thermometer": NewThermometer,
			"switch":      NewSwitch,
		},
		outputFactories: map[string]OutputDeviceFactory{
			"gpio":     NewGPIO,
			"heater":   NewHeater,
			"cooler":   NewCooler,
			"recorder": NewRecorder,
		},
	}
	return a
}

//Each input and output device has a config method that returns a Pin with
//the required fields poplulated with helpful values.
func GetTypes() map[string]ConfigHelper{
	t := Thermometer{}
	s := Switch{}
	g := GPIO{}
	h := Heater{}
	c := Cooler{}
	r := Recorder{}
	return map[string]ConfigHelper{
		"thermometer": t.Config(),
		"switch":      s.Config(),
		"gpio":     g.Config(),
		"heater":   h.Config(),
		"cooler":   c.Config(),
		"recorder": r.Config(),
	}
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
