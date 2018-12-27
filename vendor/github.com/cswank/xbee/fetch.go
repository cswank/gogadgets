package xbee

import (
	"go.bug.st/serial.v1"
)

func ReadMessage(port serial.Port, opts ...func()) (Message, error) {
	var msg Message

	for {
		getDelimiter(port)
		d := make([]byte, 2)
		n, err := port.Read(d)
		if err != nil || n != 2 {
			return msg, err
		}

		l, err := GetLength(d)
		if err != nil {
			return msg, err
		}

		d = []byte{}
		d, err = getBody(d, port, int(l+1))
		if err != nil {
			return msg, err
		}

		msg, err := NewMessage(d)
		if err == nil {
			return msg, nil
		}
	}
}

func getDelimiter(port serial.Port) {
	for {
		d := make([]byte, 1)
		n, err := port.Read(d)
		if err != nil || n != 1 {
			continue
		}

		if d[0] == 0x7E {
			return
		}
	}
}

func getBody(data []byte, port serial.Port, l int) ([]byte, error) {
	d := make([]byte, l)
	n, err := port.Read(d)
	if err != nil {
		return nil, err
	}
	d = append(data, d[:n]...)
	if n == l {
		return d, nil
	}
	return getBody(d, port, l-n)
}
