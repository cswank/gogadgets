/*
       board (master)         board 2 
   sub  pub  reply        sub  pub  reply
    |    |                 |    |
    |    -------------------    |
    -----------------------------
*/

package gogadgets

import (
	"fmt"
	"log"
	"encoding/json"
	"github.com/vaughan0/go-zmq"
)

type Sockets struct {
	masterHost string
	pubPort int
	subPort int
	isMaster bool
	ctx *zmq.Context
	sub *zmq.Socket
	subChan *zmq.Channels
	pub *zmq.Socket
	pubChan *zmq.Channels
}

func (s *Sockets) GetUID() string {
	return "zmq sockets"
}

func (s *Sockets) Start(in <-chan Message, out chan<- Message) {
	err := s.getSockets()
	defer s.ctx.Close()
	defer s.sub.Close()
	defer s.pub.Close()
	defer s.subChan.Close()
	defer s.pubChan.Close()
	if err != nil {
		log.Println("zmq sockets had a problem", err)
	}
	keepGoing := true
	for keepGoing {
		select {
		case data := <-s.subChan.In():
			msg := &Message{}
			json.Unmarshal(data[1], msg)
			msg.Sender = "zmq sockets"
			out<- *msg
		case msg := <-in:
			if msg.Type == COMMAND && msg.Body == "shutdown" {
				keepGoing = false
			}
			b, err := json.Marshal(msg)
			if err != nil {
				log.Println("zmq sockets had a problem", err)
			} else {
				s.pubChan.Out()<- [][]byte{
					[]byte(msg.Type),
					b,
				}
			}
		case err = <-s.subChan.Errors():
			log.Println(err)
		}
	}
}

func (s *Sockets) getSockets() (err error) {
	if s.masterHost == "localhost" || s.masterHost == "" {
		s.isMaster = true
		err = s.getMasterSockets()
	} else {
		err = s.getClientSockets()
	}
	return err
}

func (s *Sockets) getMasterSockets() (err error) {
	s.ctx, err = zmq.NewContext()
	if err != nil {
		return err
	}
	
	s.pub, err = s.ctx.Socket(zmq.Pub)
	if err != nil {
		return err
	}
	if err = s.pub.Bind(fmt.Sprintf("tcp://*:%d", s.pubPort)); err != nil {
		return err
	}
	s.pubChan = s.pub.Channels()

	sub, err := s.ctx.Socket(zmq.Sub)
	if err != nil {
		return err
	}
	s.sub = sub
	if err = s.sub.Bind(fmt.Sprintf("tcp://*:%d", s.subPort)); err != nil {
		return err
	}
	s.sub.Subscribe([]byte(""))
	s.subChan = s.sub.Channels()
	return err
}

func (s *Sockets) getClientSockets() (err error) {
	s.ctx, err = zmq.NewContext()
	if err != nil {
		return err
	}
	
	s.pub, err = s.ctx.Socket(zmq.Pub)
	if err != nil {
		return err
	}
	if err = s.pub.Connect(fmt.Sprintf("tcp://*:%d", s.masterHost, s.pubPort)); err != nil {
		return err
	}

	sub, err := s.ctx.Socket(zmq.Sub)
	
	if err != nil {
		return err
	}
	if err = sub.Connect(fmt.Sprintf("tcp://*:%d", s.masterHost, s.subPort)); err != nil {
		return err
	}
	sub.Subscribe([]byte(""))
	s.subChan = sub.Channels()
	s.sub = sub
	return err
}
