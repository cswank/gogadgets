package gogadgets

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"

	serial "go.bug.st/serial.v1"
)

/*
XBee reads adc remote xbees.
*/

type PortFactory func(string, *serial.Mode) (serial.Port, error)

var (
	serialFactory PortFactory
)

func Moisture(v float64) (float64, string) {
	return v / 1000.0, "%"
}

func TMP36(v float64) (float64, string) {
	return v, "F"
}

type xbeeHeader struct {
	Type            uint8
	Addr            uint64
	ShortAddr       uint16
	Opts            uint8
	Samples         uint8
	DigitalChanMask uint16
	AnalogChanMask  uint8
}

type XBeeMsg struct {
	header  xbeeHeader
	payload []byte
}

func NewXBeeMsg(d []byte) (XBeeMsg, error) {
	var h xbeeHeader
	buf := bytes.NewReader(d)
	err := binary.Read(buf, binary.BigEndian, &h)
	var msg XBeeMsg
	if err != nil {
		return msg, err
	}
	msg = XBeeMsg{
		header:  h,
		payload: d[16:],
	}
	return msg, msg.Check()
}

func (x *XBeeMsg) GetAnalog() []float64 {
	var d []byte
	if x.header.DigitalChanMask == 0 {
		d = x.payload[:len(x.payload)-1]
	} else {
		d = x.payload[2 : len(x.payload)-1]
	}
	out := make([]float64, len(d)/2)
	for i := 0; i < len(out); i++ {
		buf := bytes.NewReader(d[i*2 : i*2+2])
		var x uint16
		err := binary.Read(buf, binary.BigEndian, &x)
		if err != nil {
			log.Println("counldn't parse analog data")
		}
		out[i] = 1200.0 * float64(x) / float64(1023)
	}
	return out
}

func (x *XBeeMsg) Check() error {
	//cs := x.Payload[len(x.Payload)-1:]
	return nil
}

func (x *XBeeMsg) GetAddr() string {
	return "theaddr"
}

type XBee struct {
	port serial.Port
	conv map[string][]func(float64) (float64, string)
}

type XBeeConfig map[string][]string

func (x XBeeConfig) getConv() map[string][]func(float64) (float64, string) {
	out := map[string][]func(float64) (float64, string){}
	for k, v := range x {
		out[k] = []func(float64) (float64, string){}
		for _, n := range v {
			switch n {
			case "moisture":
				out[k] = append(out[k], Moisture)
			case "tmp36":
				out[k] = append(out[k], TMP36)
			}
		}
	}
	return out
}

func NewXBee(pin *Pin) (InputDevice, error) {
	p, ok := pin.Args["port"].(string)
	if !ok {
		return nil, fmt.Errorf(`unable to create serial port for XBee, pin.Args["port"] should be the path to a serial device`)
	}

	conv, ok := pin.Args["convert"].(XBeeConfig)
	if !ok {
		return nil, fmt.Errorf(`can't create xbee, no convertions`)
	}

	mode := &serial.Mode{}

	port, err := serialFactory(p, mode)
	if err != nil {
		return nil, fmt.Errorf(`unable to create serial port for XBee, err: %v`, err)
	}

	return &XBee{
		port: port,
		conv: conv.getConv(),
	}, nil
}

func (x *XBee) Start(ch <-chan Message, val chan<- Value) {
	go x.listen(val)
	for {
		<-ch
	}
}

func (x *XBee) listen(ch chan<- Value) {
	msgCh := make(chan XBeeMsg)
	go x.readPackets(msgCh)
	for {
		msg := <-msgCh
		c, ok := x.conv[msg.GetAddr()]
		if !ok {
			log.Println("ignoring message from unknown xbee:", msg.GetAddr())
			continue
		}
		for i, v := range msg.GetAnalog() {
			f, u := c[i](v)
			ch <- Value{
				Value: f,
				Units: u,
			}
		}
	}
}

func (x *XBee) getAnalog(msg XBeeMsg, data []byte) []uint16 {
	out := make([]uint16, len(data[:len(data)])/2)
	for i, d := range out {
		buf := bytes.NewReader(data[i*2 : (i*2)+2])
		err := binary.Read(buf, binary.LittleEndian, &d)
		if err != nil {
			log.Println(err)
		}
		out[i] = d
	}
	return out
}

func (x *XBee) readPackets(ch chan XBeeMsg) {
	for {
		b := x.readByte()
		if b != 0x7E {
			continue
		}
		buf := make([]byte, 2)
		n, err := x.port.Read(buf)
		if err != nil || n != 2 {
			continue
		}
		l := x.getLen(buf)
		buf = make([]byte, l)
		n, err = x.port.Read(buf)
		if err != nil || n != l {
			continue
		}
		msg, err := NewXBeeMsg(buf)
		if err != nil {
			continue
		}
		ch <- msg
	}
}

func (x *XBee) getLen(d []byte) int {
	var l int16
	buf := bytes.NewReader(d)
	err := binary.Read(buf, binary.BigEndian, &l)
	if err != nil {
		log.Println(err)
	}
	return int(l) + 1
}

func (x *XBee) readByte() byte {
	buf := make([]byte, 1)
	for {
		n, err := x.port.Read(buf)
		if err == nil && n == 1 {
			break
		}
	}
	return buf[0]
}

func (x *XBee) GetValue() *Value {
	return &Value{}
}

func (x *XBee) Config() ConfigHelper {
	return ConfigHelper{}
}
