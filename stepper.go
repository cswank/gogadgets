package gogadgets

type Stepper struct {
	ms1    OutputDevice
	ms2    OutputDevice
	ms3    OutputDevice
	pwm    OutputDevice
	dir    OutputDevice
	status bool
}

func NewStepper(pin *Pin) (OutputDevice, error) {
	gpios := make([]*GPIO, 3)
	for i, k := range []string{"ms1", "ms2", "ms3", "dir"} {
		p := pin.Pins[k]
		ms, err := newGPIO(&p)
		if err != nil {
			return nil, err
		}
		gpios[i] = ms
	}

	p := pin.Pins["pwm"]
	pwm, err := NewPWM(&p)
	if err != nil {
		return nil, err
	}
	return &Stepper{
		ms1: gpios[0],
		ms2: gpios[1],
		ms3: gpios[2],
		dir: gpios[3],
		pwm: pwm,
	}, nil
}

func (s *Stepper) Commands(location, name string) *Commands {
	return nil
}

func (s *Stepper) Update(msg *Message) bool {
	return false
}

func (s *Stepper) On(val *Value) error {
	if val == nil {
		val = &Value{Value: 100.0, Units: "rpm"}
	}
	v, ok := val.Value.(float64)
	if !ok {
		return nil
	}

	if v < 0.0 {
		s.dir.On(nil)
	} else if v > 0.0 {
		s.dir.Off()
	}

	var hz float64
	switch val.Units {
	case "steps":
	case "rpm":
	}

	return s.pwm.On(&Value{Value: hz})
}

func (s *Stepper) Status() map[string]bool {
	return map[string]bool{
		"ms1": s.ms1.Status()["gpio"],
		"ms2": s.ms2.Status()["gpio"],
		"ms3": s.ms3.Status()["gpio"],
		"dir": s.dir.Status()["gpio"],
	}
}

func (s *Stepper) Off() error {
	s.pwm.Off()
	return nil
}
