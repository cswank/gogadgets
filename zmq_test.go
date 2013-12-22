package gogadgets

import (
	"fmt"
	"time"
	"testing"
	"encoding/json"
	"github.com/vaughan0/go-zmq"
)

func TestSockets(t *testing.T) {
	s := Sockets{}
	input := make(chan Message)
	output := make(chan Message)
	go s.Start(input, output)
	
	ctx, err := zmq.NewContext()
	defer ctx.Close()
	if err != nil {
		t.Fatal(err)
	}
	
	pub, err := ctx.Socket(zmq.Pub)
	defer pub.Close()
	if err != nil {
		t.Fatal(err)
	}
	err = pub.Connect("tcp://localhost:6112")
	if err != nil {
		t.Fatal(err)
	}

	// sub, err := ctx.Socket(zmq.Pub)
	// defer sub.Close()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// sub.Connect("tcp://localhost:6111"); err != nil {
	// 	t.Fatal(err)
	// }
	msg := Message{
		Type: "command",
		Body: "testing testing",
	}
	time.Sleep(10 * time.Millisecond)
	fmt.Println(msg)
	b, err := json.Marshal(msg)
	if err != nil {
		t.Fatal(err)
	}
	pub.Send(
		[][]byte{
			[]byte(msg.Type),
			b,
		},
	)
	time.Sleep(10 * time.Millisecond)
}
