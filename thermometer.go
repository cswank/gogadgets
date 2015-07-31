package gogadgets

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cswank/gogadgets/utils"
)

/*Reads temperature from a Dallas 1-Wire thermometer and
sends that temperature to the rest of the system.

on ubuntu install dtc (patched)
  wget -c https://raw.github.com/RobertCNelson/tools/master/pkgs/dtc.sh
  chmod +x dtc.sh
  ./dtc.sh

  echo '
/dts-v1/;
/plugin/;

/ {
compatible = "ti,beaglebone", "ti,beaglebone-black";

part-number = "BB-W1";
version = "00A0";

exclusive-use =
"P9.22",
"gpio0_2";

fragment@0 {
               target = <&am33xx_pinmux>;
               __overlay__ {
dallas_w1_pins: pinmux_dallas_w1_pins {
pinctrl-single,pins = < 0x150 0x37 >;
};
               };
};

fragment@1 {
               target = <&ocp>;
               __overlay__ {
       onewire@0 {
       compatible      = "w1-gpio";
       pinctrl-names   = "default";
       pinctrl-0       = <&dallas_w1_pins>;
       status          = "okay";

       gpios = <&gpio1 2 0>;
       };
         };
};
};' > BB-W1-00A0.dts

  dtc -O dtb -o BB-W1-00A0.dtbo -b 0 -@ BB-W1-00A0.dts
  cp BB-W1-00A0.dtbo /lib/firmware/
  echo BB-W1:00A0 > /sys/devices/bone_capemgr.9/slots

*/

type Thermometer struct {
	devicePath string
	units      string
	value      float64
	sleep      time.Duration
	lock       sync.Mutex
}

func NewThermometer(pin *Pin) (InputDevice, error) {
	var therm *Thermometer
	var err error
	if pin.OneWirePath == "" {
		pin.OneWirePath = "/sys/bus/w1/devices/%s/w1_slave"
	}
	if pin.Sleep == 0 {
		pin.Sleep = 5 * time.Second
	}
	path := fmt.Sprintf(pin.OneWirePath, pin.OneWireId)
	if pin.OneWireId == "" || !utils.FileExists(path) {
		err = errors.New(fmt.Sprintf("invalid one-wire device path: %s", pin.OneWireId))
		return therm, err
	}
	therm = &Thermometer{
		devicePath: path,
		units:      pin.Units,
		sleep:      pin.Sleep,
		lock:       pin.Lock,
	}
	return therm, err
}

func (t *Thermometer) Config() ConfigHelper {
	return ConfigHelper{
		PinType: "thermometer",
		Fields: map[string][]string{
			"oneWireId": []string{},
		},
		Units: []string{"C", "F"},
	}
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
			t.value = val.Value.(float64)
			out <- *val
		}
		time.Sleep(t.sleep)
	}
}

//The 1-wire craps out once in a while and a value less than zero is a sign
//that something went wrong.  Ususally the subsequent temperature value
//is valid.
func (t *Thermometer) isValid(value, previous *Value) bool {
	return previous == nil || math.Abs(previous.Value.(float64)-value.Value.(float64)) < 10.0
}

// Linux (with certain kernels) on a Beaglebone and Raspberry Pi have
// a file based interface to the Dallas 1-wire devices.  This reads
// from that file.
func (t *Thermometer) readFile() (v *Value, err error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	b, err := ioutil.ReadFile(t.devicePath)
	if err != nil {
		return v, err
	}
	return t.parseValue(string(b))
}

//arseValue gets the actual tempreature from the 1-wire interface
//sysfs file.
func (t *Thermometer) parseValue(val string) (*Value, error) {
	var v *Value
	start := strings.Index(val, "t=")
	if start == -1 {
		return v, errors.New("could not parse temp")
	}
	temperatureStr := val[start+2:]
	temperatureStr = strings.Trim(temperatureStr, "\n")
	temperature, err := strconv.ParseFloat(temperatureStr, 64)
	if err != nil {
		return v, err
	}
	value := temperature / 1000.0
	if t.units == "F" || t.units == "f" {
		value = value*1.8 + 32.0
	}
	return &Value{
		Value: value,
		Units: t.units,
	}, nil
}

//This is an InputDevice, so it must have a Start.
func (t *Thermometer) Start(in <-chan Message, out chan<- Value) {
	temperature := make(chan Value)

	e := make(chan error)
	go t.getTemperature(temperature, e)
	for {
		select {
		case <-in:
			// do nothing
		case val := <-temperature:
			out <- val
		case err := <-e:
			log.Println(fmt.Sprintf("error reading thermometer %s", t.devicePath), err)
		}
	}
}
