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
	host string
	pubPort int
	subPort int
	isMaster bool
	ctx *zmq.Context
	sub *zmq.Socket
	subChan *zmq.Channels
	pub *zmq.Socket
	pubChan *zmq.Channels
}

func NewClientSockets(host string) (*Sockets, error) {
	s := &Sockets{
		host: host,
		pubPort: 6112,
		subPort: 6111,
	}
	err := s.getClientSockets()
	return s, err
}

func (s *Sockets) GetUID() string {
	return "zmq sockets"
}

func (s *Sockets) Send(cmd string) {
	msg := &Message{
		Type: COMMAND,
		Body: cmd,
	}
	b, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("zmq sockets had a problem", err)
	} else {
		s.pub.Send([][]byte{
			[]byte(msg.Type),
			b,
		})
	}
}

func (s *Sockets) Start(in <-chan Message, out chan<- Message) {
	err := s.getSockets()
	defer s.Close()
	if err != nil {
		log.Println("zmq sockets had a problem", err)
	}
	keepGoing := true
	for keepGoing {
		select {
		case data := <-s.subChan.In():
			if len(data) == 2 {
				msg := &Message{}
				json.Unmarshal(data[1], msg)
				msg.Sender = "zmq sockets"
				out<- *msg
			} else {
				log.Println("zmq received an improper message", data)
			}
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

func (s *Sockets) Close() {
	s.ctx.Close()
	s.sub.Close()
	s.pub.Close()
	s.subChan.Close()
	s.pubChan.Close()
}

func (s *Sockets) getSockets() (err error) {
	if s.host == "localhost" || s.host == "" {
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
	if err = s.pub.Connect(fmt.Sprintf("tcp://%s:%d", s.host, s.subPort)); err != nil {
		return err
	}
	
	s.sub, err = s.ctx.Socket(zmq.Sub)
	
	if err != nil {
		return err
	}
	if err = s.sub.Connect(fmt.Sprintf("tcp://%s:%d", s.host, s.pubPort)); err != nil {
		return err
	}
	s.sub.Subscribe([]byte(""))
	return err
}
