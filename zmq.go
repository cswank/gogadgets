package gogadgets

import (
	"fmt"
	"log"
	"encoding/json"
	"github.com/vaughan0/go-zmq"
)

type Sockets struct {
	GoGadget
	ctx *zmq.Context
	reply *zmq.Channels
	sub *zmq.Channels
	pub *zmq.Socket
}

func (s *Sockets) Start(in <-chan Message, out chan<- Message) {
	err := s.getSockets()
	defer s.ctx.Close()
	defer s.sub.Close()
	defer s.pub.Close()
	defer s.reply.Close()
	if err != nil {
		log.Println("zmq sockets had a problem", err)
	}
	fmt.Println("select")
	select {
	case data := <-s.sub.In():
		msg := &Message{}
		json.Unmarshal(data[1], msg)
		fmt.Println("sockets got", msg)
		out<- *msg
	case data := <-s.reply.In():
		msg := &Message{}
		json.Unmarshal(data[1], msg)
		out<- *msg
	case msg := <-in:
		b, err := json.Marshal(msg)
		if err != nil {
			log.Println("zmq sockets had a problem", err)
		} else {
			s.pub.Send([][]byte{
				[]byte(msg.Type),
				b,
			})
		}
	case err = <-s.sub.Errors():
		log.Println(err)
	}
}

func (s *Sockets) getSockets() (err error) {
	s.ctx, err = zmq.NewContext()
	if err != nil {
		return err
	}
	reply, err := s.ctx.Socket(zmq.Rep)
	if err != nil {
		return err
	}
	if err = reply.Bind("tcp://*:6113"); err != nil {
		return err
	}
	s.reply = reply.Channels()
	
	s.pub, err = s.ctx.Socket(zmq.Pub)
	if err != nil {
		return err
	}
	if err = s.pub.Bind("tcp://*:6111"); err != nil {
		return err
	}

	sub, err := s.ctx.Socket(zmq.Sub)
	
	if err != nil {
		return err
	}
	if err = sub.Bind("tcp://*:6112"); err != nil {
		return err
	}
	sub.Subscribe([]byte(""))
	s.sub = sub.Channels()
	return err
}
