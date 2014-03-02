package gogadgets

import (
	"fmt"
	"time"
	"log"
	"errors"
	"io/ioutil"
	"strings"
	"strconv"
	"bitbucket.org/cswank/gogadgets/utils"
)

type Thermometer struct {
	InputDevice
	devicePath string
	units string
	value float64
}

func NewThermometer(pin *Pin) (InputDevice, error) {
	var therm *Thermometer
	var err error
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

func (t *Thermometer) GetValue() *Value {
	return &Value{
		Value: t.value,
		Units: t.units,
	}
}

func (t *Thermometer) getTemperature(out chan Value, err chan error) {
	var previousTemperature *Value
	for {
		val, e := t.readFile()
		if e == nil && t.isValid(val, previousTemperature) {
			previousTemperature = val
			out<- *val
		} else {
			err<- e
		}
		time.Sleep(5 * time.Second)
	}
}

func (t *Thermometer) isValid(value, previous *Value) (isValid bool) {
	if previous == nil {
		isValid = true
	} else if value.Value.(float64) < 0.0 { //the 1-wire craps out once in a while and a value less than zero is a sign
		isValid = false
	} else {
		isValid = true
	}
	return isValid
}

func (t *Thermometer) readFile() (v *Value, err error) {
	b, err := ioutil.ReadFile(t.devicePath)
	if err != nil {
		return v, err
	}
	return t.parseValue(string(b))
}

func (t *Thermometer) parseValue(val string) (v *Value, err error) {
	start := strings.Index(val, "t=")
	if start == -1 {
		return v, errors.New("could not parse temp")
	}
	temperatureStr := val[start + 2:]
	temperatureStr = strings.Trim(temperatureStr, "\n")
	temperature, err := strconv.ParseFloat(temperatureStr, 64)
	if err == nil {
		t.value = temperature / 1000.0
		v = &Value{
			Value: t.value,
			Units: t.units,
		}
	}
	return v, err
}

func (t *Thermometer) Start(in <-chan Message, out chan<- Value) {
	temperature := make(chan Value)
	e := make(chan error)
	go t.getTemperature(temperature, e)
	for {
		select {
		case <-in:
			// do nothing
		case val := <-temperature:
			out<- val
		case err := <-e:
			log.Println(fmt.Sprintf("error reading thermometer %s", t.devicePath), err)
		}
	}
}
