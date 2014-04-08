package gogadgets

import (
	"time"
	"errors"
	"path"
	"filepath"
)


const (
	DEVPATH = "/sys/devices/ocp.*/pwm_test_P%s_%s.*"
	TREEPATH = "/sys/devices/bone_capemgr.*/slots"
)


// echo am33xx_pwm > /sys/devices/bone_capemgr.9/slots
// echo bone_pwm_P8_13 > /sys/devices/bone_capemgr.9/slots
// /sys/devices/ocp.3/pwm_test_P8_13.15
type PWM struct {
	period int
}

func NewPWM(pin *Pin) (OutputDevice, error) {
	if !ok {
		return errors.New("invalid config for pwm")
	}
	err := writePWMDeviceTree(pin.Port, pin.Pin)
	if err != nil {
		return nil, err
	}
	period, err := setupPWM(pin)
	pwm := &PWM{
		period: period,
	}
	return pwm, nil
}

func (p *PWM) Update(msg *Message) {
	
}

func (p *PWM) On(val *Value) error {
	return nil
}

func (p *PWM) Status() interface{} {
	return false
}

func (p *PWM) Off() error {
	return nil
}

func setupPWM(pin *Pin) error {
	g := filepath.Glob(fmt.Sprintf(DEVPATH, port, pin))
	if len(g) != 1 {
		return 0, errors.New(fmt.Sprintf("couldn't find device path for PWM port %s pin %s", pin.Port, pin.Pin))
	}
	p := g[0]
	period := int(pin.Frequency / time.Nanosecond)
	path.Join(p, "period")
	err := ioutil.WriteFile(p, []byte(fmt.Sprintf("%d", ns), os.ModeDevice))
	return period, err
}

func writePWMDeviceTree(port, pin string) error {
	treePath, err := getTreePath()
	if err != nil {
		return err
	}
	port, ok := Pins["pwm"][port]
	if !ok {
		return errors.New(fmt.Sprintf("invalid port: %s", port))
	}
	val, ok := port[pin]
	if !ok {
		return errors.New(fmt.Sprintf("invalid pin: %s", pin))
	}
	err = ioutil.WriteFile(treePath, []byte("am33xx_pwm"), os.ModeDevice)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(treePath, []byte(val), os.ModeDevice)
}

func getTreePath() (string, error) {
	g := filepath.Glob(TREEPATH)
	if len(g) != 1 {
		return "", errors.New("couldn't find device tree path for slots")
	}
	return g[0], nil
}
