package xbee

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
)

type header struct {
	Type            uint8
	Addr            uint64
	ShortAddr       uint16
	Opts            uint8
	Samples         uint8
	DigitalChanMask uint16
	AnalogChanMask  uint8
}

type Message struct {
	header
	frame []byte
}

var (
	adc = []string{"adc0", "adc1", "adc2", "adc3"}
	dio = []string{"dio0", "dio1", "dio2", "dio3", "dio4", "dio5", "dio6", "dio7", "dio8", "dio9", "dio10", "dio11", "dio12"}
)

func GetLength(data []byte) (uint16, error) {
	buf := bytes.NewReader(data)
	var l uint16
	return l, binary.Read(buf, binary.BigEndian, &l)
}

//NewMessage parses a xbee messsage (from the 3rd byte on).
//The first 3 bytes are used to detect the beginning of a new
//message and the length of the message so are only useful for
//ingesting data from an xbee.
func NewMessage(data []byte) (Message, error) {
	var h header
	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.BigEndian, &h)
	var msg Message
	if err != nil {
		return msg, err
	}
	msg = Message{
		header: h,
		frame:  data,
	}
	if !msg.Check() {
		err = fmt.Errorf("message failed checksum: 0x%x", msg.getChecksum())
	}
	return msg, err
}

func (x *Message) payload() []byte {
	return x.frame[16 : len(x.frame)-1]
}

//GetAnalog returns a map of the adc pins and their voltages (in mV).
//If no pins are configured for adc then it returns an empty map.
func (x *Message) GetAnalog() (map[string]float64, error) {
	m := map[string]float64{}
	if x.AnalogChanMask == 0 {
		return m, nil
	}

	vals, err := x.getRawAnalog()
	if err != nil {
		return m, err
	}

	var j int
	for i, o := range adc {
		if uint8(1<<uint8(i))&x.AnalogChanMask > 0 {
			x := vals[j]
			m[o] = 1200.0 * float64(x) / 1023.0
			j++
		}
	}
	return m, nil
}

func (x *Message) getRawAnalog() ([]uint16, error) {
	var d []byte
	payload := x.payload()
	if x.DigitalChanMask == 0 {
		d = payload
	} else {
		d = payload[2:]
	}

	f := make([]uint16, len(d)/2)
	buf := bytes.NewReader(d)
	return f, binary.Read(buf, binary.BigEndian, &f)
}

func (x *Message) GetDigital() (map[string]bool, error) {
	m := map[string]bool{}
	if x.DigitalChanMask == 0 {
		return m, nil
	}

	val, err := x.getRawDigital()
	if err != nil {
		return m, err
	}

	var j int
	for i, o := range dio {
		if uint16(1<<uint16(i))&x.DigitalChanMask > 0 {
			x := val & (1 << uint16(i))
			m[o] = x > 0
			j++
		}
	}
	return m, nil
}
func (x *Message) getRawDigital() (uint16, error) {
	payload := x.payload()[:3]
	var val uint16
	buf := bytes.NewReader(payload)
	return val, binary.Read(buf, binary.BigEndian, &val)
}

func (x *Message) getChecksum() byte {
	var total byte
	for _, item := range x.frame[:len(x.frame)-1] {
		total += item
	}
	return 0xff - total
}

func (x *Message) Check() bool {

	cs := x.frame[len(x.frame)-1]
	return x.getChecksum() == cs
}

func (x *Message) GetAddr() string {
	return strconv.FormatUint(x.Addr, 16)
}
