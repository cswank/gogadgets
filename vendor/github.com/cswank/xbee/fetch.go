package xbee

import (
	"log"

	"go.bug.st/serial.v1"
)

func ReadMessage(port serial.Port) Message {
	for {
		getDelimiter(port)
		d := make([]byte, 2)
		n, err := port.Read(d)
		if err != nil || n != 2 {
			continue
		}

		l, err := GetLength(d)
		if err != nil {
			continue

		}

		d = []byte{}
		d, err = getBody(d, port, int(l+1))
		if err != nil {
			continue
		}

		msg, err := NewMessage(d)
		if err == nil {
			return msg
		}
	}
}

func getDelimiter(port serial.Port) {
	for {
		d := make([]byte, 1)
		n, err := port.Read(d)
		if err != nil || n != 1 {
			log.Println(n, err)
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
