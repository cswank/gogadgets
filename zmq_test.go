package gogadgets

// import (
// 	"encoding/json"
// 	"fmt"
// 	"github.com/vaughan0/go-zmq"
// 	"math/rand"
// 	"testing"
// 	"time"
// )

// func TestSockets(t *testing.T) {
// 	pubPort := 1024 + rand.Intn(65535-1024)
// 	subPort := pubPort + 1
// 	s := Sockets{
// 		host:    "localhost",
// 		pubPort: pubPort,
// 		subPort: subPort,
// 	}
// 	input := make(chan Message)
// 	output := make(chan Message)
// 	go s.Start(input, output)

// 	ctx, err := zmq.NewContext()
// 	defer ctx.Close()
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	pub, err := ctx.Socket(zmq.Pub)
// 	defer pub.Close()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	err = pub.Connect(fmt.Sprintf("tcp://localhost:%d", subPort))
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	chans := pub.Channels()
// 	defer chans.Close()

// 	sub, err := ctx.Socket(zmq.Sub)
// 	defer sub.Close()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	if err = sub.Connect(fmt.Sprintf("tcp://localhost:%d", pubPort)); err != nil {
// 		t.Fatal(err)
// 	}
// 	sub.Subscribe([]byte(""))

// 	msg := Message{
// 		Type: "command",
// 		Body: "testing testing",
// 	}

// 	b, _ := json.Marshal(msg)

// 	go func() {
// 		time.Sleep(500 * time.Millisecond)
// 		chans.Out() <- [][]byte{
// 			[]byte(msg.Type),
// 			b,
// 		}
// 	}()
// 	<-output
// 	go func() {
// 		time.Sleep(50 * time.Millisecond)
// 		input <- msg
// 	}()
// 	parts, err := sub.Recv()

// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if len(parts) != 2 {
// 		t.Error(len(parts))
// 	}
// 	if string(parts[0]) != "command" {
// 		t.Error(string(parts[0]))
// 	}

// 	msg = Message{
// 		Type: "command",
// 		Body: "shutdown",
// 	}
// 	input <- msg
// }
