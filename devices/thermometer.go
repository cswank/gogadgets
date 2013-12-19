package devices

import (
	"fmt"
	"time"
	"log"
	"errors"
	"io/ioutil"
	"strings"
	"strconv"
	"bitbucket.com/cswank/gogadgets/utils"
	"bitbucket.com/cswank/gogadgets/models"
)

type Thermometer struct {
	InputDevice
	devicePath string
	units string
}

func NewThermometer(pin *models.Pin) (therm *Thermometer, err error) {
	path := fmt.Sprintf("/sys/bus/w1/devices/%s/w1_slave", pin.OneWireId)
	if pin.OneWireId == "" || !utils.FileExists(path) {
		err = errors.New(fmt.Sprintf("invalid one-wire device path: %s", pin.OneWireId))
		return therm, err
	}
	therm = &Thermometer{
		devicePath: path,
		units: "C",
	}
	return therm, err
}

func (t *Thermometer) getValue() (v *models.Value, err error) {
	b, err := ioutil.ReadFile(t.devicePath)
	if err == nil {
		return t.parseValue(string(b))
	}
	return v, err
}

func (t *Thermometer) parseValue(val string) (v *models.Value, err error) {
	start := strings.Index(val, "t=")
	if start == -1 {
		return v, errors.New("could not parse temp")
	}
	temperatureStr := val[start + 2:]
	temperatureStr = strings.Trim(temperatureStr, "\n")
	temperature, err := strconv.ParseFloat(temperatureStr, 64)
	if err == nil {
		v = &models.Value{
			Value: temperature / 1000.0,
			Units: t.units,
		}
	}
	return v, err
}

func (t *Thermometer) Start(stop <-chan bool, out chan<- models.Value) {
	for {
		select {
		case <-stop:
			return
		case <-time.After(5 * time.Second):
			val, err := t.getValue()
			if err == nil {
				out<- *val
			} else {
				log.Println(fmt.Sprintf("error reading thermometer %s", t.devicePath))
			}
			
		}
	}
}
