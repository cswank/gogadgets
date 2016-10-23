package xbee

import (
	"log"

	"go.bug.st/serial.v1"
)

var (
	verbose bool
)

func Verbose() {
	verbose = true
}

func logger(i ...interface{}) {
	if verbose {
		log.Print(i...)
	}
}

func ReadMessage(port serial.Port, opts ...func()) Message {
	for _, o := range opts {
		o()
	}

	for {
		getDelimiter(port)
		d := make([]byte, 2)
		n, err := port.Read(d)
		if err != nil || n != 2 {
			logger(err, n)
			continue
		}

		l, err := GetLength(d)
		if err != nil {
			logger(err)
			continue

		}

		d = []byte{}
		d, err = getBody(d, port, int(l+1))
		if err != nil {
			logger(string(d), err)
			continue
		}

		msg, err := NewMessage(d)
		if err == nil {
			logger(msg, err)
			return msg
		}
	}
}

func getDelimiter(port serial.Port) {
	for {
		d := make([]byte, 1)
		n, err := port.Read(d)
		if err != nil || n != 1 {
			logger(n, err)
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
