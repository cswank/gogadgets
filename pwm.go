package gogadgets

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
	"path/filepath"
)

const (
	NANO = 1000000000.0
)

var (
	pwmMode     = os.ModeDevice
	PWM_DEVPATH = "/sys/devices/ocp.*/pwm_test_P%s_%s.*"
	TREEPATH    = "/sys/devices/bone_capemgr.*/slots"
)

// echo am33xx_pwm > /sys/devices/bone_capemgr.9/slots
// echo bone_pwm_P8_13 > /sys/devices/bone_capemgr.9/slots
// /sys/devices/ocp.3/pwm_test_P8_13.15
type PWM struct {
	period     int
	duty       []byte
	status     bool
	runPath    string
	dutyPath   string
	periodPath string
}

func NewPWM(pin *Pin) (OutputDevice, error) {
	// err := writePWMDeviceTree(pin.Port, pin.Pin)
	// if err != nil {
	// 	return nil, err
	// }
	devPath, period, err := setupPWM(pin)
	pwm := &PWM{
		period:     period,
		duty:       []byte(fmt.Sprintf("%d", period)),
		runPath:    path.Join(devPath, "run"),
		dutyPath:   path.Join(devPath, "duty"),
		periodPath: path.Join(devPath, "period"),
	}
	return pwm, err
}

func (p *PWM) Config() ConfigHelper {
	return ConfigHelper{
		PinType: "pwm",
		Fields: map[string][]string{
			"frequency": []string{},
		},
		Pins: Pins["pwm"],
	}
}

func (p *PWM) Update(msg *Message) {

}

func (p *PWM) On(val *Value) error {
	if val != nil && val.Units == "%" {
		ioutil.WriteFile(p.runPath, []byte("0"), pwmMode)
		p.duty = p.getDuty(val.Value)
		ioutil.WriteFile(p.dutyPath, p.duty, pwmMode)
	} else {
		ioutil.WriteFile(p.runPath, []byte("0"), pwmMode)
		ioutil.WriteFile(p.dutyPath, []byte(fmt.Sprintf("%d", p.period)), pwmMode)
	}
	return ioutil.WriteFile(p.runPath, []byte("1"), pwmMode)
}

func (p *PWM) Off() error {
	ioutil.WriteFile(p.dutyPath, []byte("0"), pwmMode)
	return ioutil.WriteFile(p.runPath, []byte("0"), pwmMode)
}

func (p *PWM) Status() interface{} {
	return p.status
}

func (p *PWM) getDuty(val interface{}) []byte {
	d, ok := val.(float64)
	if !ok {
		return []byte("0")
	}
	d = math.Abs(d)
	f := (d / 100.0) * float64(p.period)
	return []byte(fmt.Sprintf("%d", int(f)))
}

func setupPWM(pin *Pin) (devPath string, period int, err error) {
	g, e := filepath.Glob(fmt.Sprintf(PWM_DEVPATH, pin.Port, pin.Pin))
	if e != nil {
		return devPath, period, e
	}
	if len(g) != 1 {
		return devPath, period, errors.New(fmt.Sprintf("couldn't find device path for PWM port %s pin %s path %s", pin.Port, pin.Pin, PWM_DEVPATH))
	}
	devPath = g[0]
	period = int(NANO / float32(pin.Frequency))
	p := path.Join(devPath, "period")
	err = ioutil.WriteFile(p, []byte(fmt.Sprintf("%d", period)), pwmMode)
	if err != nil {
		return devPath, period, err
	}
	p = path.Join(devPath, "duty")
	err = ioutil.WriteFile(p, []byte(fmt.Sprintf("%d", period)), pwmMode)
	if err != nil {
		return devPath, period, err
	}
	p = path.Join(devPath, "polarity")
	err = ioutil.WriteFile(p, []byte("0"), pwmMode)
	if err != nil {
		return devPath, period, err
	}
	return devPath, period, err
}

func writePWMDeviceTree(port, pin string) error {
	treePath, err := getTreePath()
	if err != nil {
		return err
	}
	pwm := Pins["pwm"]
	p, ok := pwm[port]
	if !ok {
		return errors.New(fmt.Sprintf("invalid port: %s", p))
	}
	val, ok := p[pin]
	if !ok {
		return errors.New(fmt.Sprintf("invalid pin: %s", pin))
	}
	err = ioutil.WriteFile(treePath, []byte("am33xx_pwm"), pwmMode)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(treePath, []byte(val), pwmMode)
}

func getTreePath() (string, error) {
	g, err := filepath.Glob(TREEPATH)
	if err != nil {
		return "", err
	}
	if len(g) != 1 {
		return "", errors.New("couldn't find device tree path for slots")
	}
	return g[0], nil
}
