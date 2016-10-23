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

type SerialFactory func(string, *serial.Mode) (serial.Port, error)

var (
	serialFactory SerialFactory
)

func Moisture(location string) func(float64) (float64, string, string, string) {
	return func(v float64) (float64, string, string, string) {
		return v / 1000.0, "%", location, "moisture"
	}
}

func TMP36(location string) func(float64) (float64, string, string, string) {
	return func(v float64) (float64, string, string, string) {
		c := (v - 0.5) * 100.0
		f := (c * 9.0 / 5.0) + 32.0
		return f, "F", location, "temperature"
	}
}

type XBee struct {
	port serial.Port
	//        addr       pin
	xbees map[string]map[string]func(float64) (float64, string, string, string)
}

type XBeeConfig struct {
	//           pin    type
	Pins     map[string]string `json:"pins"`
	Location string            `json:"location"`
}

func (x XBeeConfig) getConv(location string) map[string]func(float64) (float64, string, string, string) {
	out := map[string]func(float64) (float64, string, string, string){}
	for k, t := range x.Pins {
		switch t {
		case "moisture":
			out[k] = Moisture(location)
		case "tmp36":
			out[k] = TMP36(location)
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
	xbees := map[string]map[string]func(float64) (float64, string, string, string){}
	for addr, xbee := range m {
		xbees[addr] = xbee.getConv(xbee.Location)
	}

	return &XBee{
		port:  port,
		xbees: xbees,
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
		x, ok := x.xbees[msg.GetAddr()]
		if !ok {
			log.Println("ignoring message from unknown xbee:", msg.GetAddr())
			continue
		}
		a, err := msg.GetAnalog()
		if err != nil {
			continue
		}

		for k, v := range a {
			go func(k string, v float64) {
				f, ok := x[k]
				if !ok {
					log.Println("ignoring message from unknown xbee pin:", k, x)
				} else {
					val, u, location, name := f(v)
					ch <- Value{
						Value:    val,
						Units:    u,
						location: location,
						name:     name,
					}
				}
			}(k, v)
		}
	}
}

func (x *XBee) readMessage(ch chan<- xbee.Message) {
	for {
		msg := xbee.ReadMessage(x.port)
		ch <- msg
	}
}

func (x *XBee) GetValue() *Value {
	return &Value{}
}

func (x *XBee) Config() ConfigHelper {
	return ConfigHelper{}
}
