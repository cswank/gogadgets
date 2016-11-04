package gogadgets

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/cswank/xbee"

	"go.bug.st/serial.v1"
)

/*
XBee reads adc from remote xbees.
*/

type address struct {
	name     string
	location string
}

type SerialFactory func(string, *serial.Mode) (serial.Port, error)

type converter func(float64) (float64, string, address)

var (
	serialFactory SerialFactory
)

func moisture(location string) converter {
	return func(v float64) (float64, string, address) {
		return v / 1000.0, "%", address{"moisture", location}
	}
}

func tmp36(location string) converter {
	return func(v float64) (float64, string, address) {
		c := (v - 500.0) / 10.0
		f := c*1.8 + 32.0
		return f, "F", address{"temperature", location}
	}
}

type XBee struct {
	port serial.Port
	//        addr     pin
	adc map[string]map[string]converter
	dio map[string]map[string]address
}

type XBeeConfig struct {
	//           pin    type
	ADC map[string]string `json:"adc"`
	//           pin    name
	DIO      map[string]string `json:"dio"`
	Location string            `json:"location"`
}

func (x XBeeConfig) getDigital(location string) map[string]address {
	out := map[string]address{}
	for k, n := range x.DIO {
		out[k] = address{name: n, location: location}
	}
	return out
}

func (x XBeeConfig) getConversion(location string) map[string]converter {
	out := map[string]converter{}
	for k, t := range x.ADC {
		switch t {
		case "moisture":
			out[k] = moisture(location)
		case "tmp36":
			out[k] = tmp36(location)
		}
	}
	return out
}

func NewXBee(pin *Pin) (InputDevice, error) {
	p, ok := pin.Args["port"].(string)
	if !ok {
		return nil, fmt.Errorf(`unable to create serial port for XBee, pin.Args["port"] should be the path to a serial device`)
	}

	j, ok := pin.Args["xbees"].(string)
	if !ok {
		return nil, fmt.Errorf(`can't create xbee: %v`, pin.Args["xbees"])
	}

	var m map[string]XBeeConfig
	if err := json.Unmarshal([]byte(j), &m); err != nil {
		return nil, fmt.Errorf(`can't create xbee: %v`, err)
	}

	mode := &serial.Mode{}

	port, err := serialFactory(p, mode)
	if err != nil {
		return nil, fmt.Errorf(`unable to create serial port for XBee, err: %v`, err)
	}

	//           addr       pin
	adc := map[string]map[string]converter{}
	dio := map[string]map[string]address{}
	for addr, x := range m {
		adc[addr] = x.getConversion(x.Location)
		dio[addr] = x.getDigital(x.Location)
	}

	return &XBee{
		port: port,
		adc:  adc,
		dio:  dio,
	}, nil
}

func (x *XBee) Start(ch <-chan Message, val chan<- Value) {
	go x.listen(val)
	for {
		<-ch
	}
}

func (x *XBee) listen(ch chan<- Value) {
	msgCh := make(chan xbee.Message)
	go x.readMessage(msgCh)
	for {
		msg := <-msgCh

		x.getAnalog(msg, ch)
		x.getDigital(msg, ch)
	}
}

func (x *XBee) getDigital(msg xbee.Message, ch chan<- Value) {
	m, ok := x.dio[msg.GetAddr()]
	if !ok {
		log.Println("ignoring message from unknown xbee:", msg.GetAddr())
		return
	}

	d, err := msg.GetDigital()
	if err != nil {
		return
	}
	for k, v := range d {
		go func(k string, v bool) {
			loc, ok := m[k]
			if !ok {
				log.Println("ignoring message from unknown xbee pin:", k, x)
			} else {
				ch <- Value{
					Value:    v,
					location: loc.location,
					name:     loc.name,
				}
			}
		}(k, v)
	}
}

func (x *XBee) getAnalog(msg xbee.Message, ch chan<- Value) {
	m, ok := x.adc[msg.GetAddr()]
	if !ok {
		log.Println("ignoring message from unknown xbee:", msg.GetAddr())
		return
	}

	a, err := msg.GetAnalog()
	if err != nil {
		return
	}

	for k, v := range a {
		go func(k string, v float64) {
			f, ok := m[k]
			if !ok {
				log.Println("ignoring message from unknown xbee pin:", k, x)
			} else {
				val, u, loc := f(v)
				ch <- Value{
					Value:    val,
					Units:    u,
					location: loc.location,
					name:     loc.name,
				}
			}
		}(k, v)
	}
}

func (x *XBee) readMessage(ch chan<- xbee.Message) {
	for {
		msg := xbee.ReadMessage(x.port, xbee.Verbose)
		ch <- msg
	}
}

func (x *XBee) GetValue() *Value {
	return &Value{}
}

func (x *XBee) Config() ConfigHelper {
	return ConfigHelper{}
}
