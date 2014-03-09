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

//Sockets fufills the GoGadget interface and is
//added to each Gadget system by App.  It provides
//a way to connect multiple gadgets systems together
//as a single system, and also provides a way for
//an external UI to control the system.
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

//This is the main loop for Sockets.  It listens for chann Messages
//from inside the system and sends it to external listeners (like a
//UI), and listens for external messages and sends them along to
//the internal system.
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

//Two GoGadgets systems can be joined together into a single system.
//One of the GoGadgets system must be declared as being the master,
//and the zmq sockets for the master are created here.
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

//This creates the zmq sockets for a GoGadget system that is not the master
//system.
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
