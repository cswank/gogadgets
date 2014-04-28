/*
       master             board 2 (or ui)
   sub  pub  reply        sub  pub  reply
    |    |                 |    |
    |    -------------------    |
    -----------------------------
*/

package gogadgets

import (
	"encoding/json"
	"fmt"
	"github.com/vaughan0/go-zmq"
	"bitbucket.org/cswank/gogadgets/models"
	"log"
)

//Sockets fufills the GoGadget interface and is
//added to each Gadget system by App.  It provides
//a way to connect multiple gadgets systems together
//as a single system, and also provides a way for
//an external UI to control the system.
type Sockets struct {
	host    string
	pubPort int
	subPort int
	ctx     *zmq.Context
	sub     *zmq.Socket
	subChan *zmq.Channels
	pub     *zmq.Socket
	pubChan *zmq.Channels
}

func NewClientSockets(host string) (*Sockets, error) {
	s := &Sockets{
		host:    host,
		subPort: 6111,
		pubPort: 6112,
	}
	err := s.getClientSockets()
	return s, err
}

func (s *Sockets) Send(cmd string) {
	msg := &models.Message{
		Type: models.COMMAND,
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

//Sockets listens for chann Messages from inside the system and
//sends it to external listeners (like a UI), and listens for
//external messages and sends them along to the internal system.
func (s *Sockets) Start(in <-chan models.Message, out chan<- models.Message) {
	err := s.getSockets()
	defer s.Close()
	if err != nil {
		log.Println("zmq sockets had a problem", err)
	}
	for {
		select {
		case data := <-s.subChan.In():
			s.sendMessageIn(data, out)
		case msg := <-in:
			s.sendMessageOut(msg)
		case err = <-s.subChan.Errors():
			log.Println(err)
		}
	}
}

//A message that came from inside this gogadgets system
//is sent to outside clients (ui, connected gogadget systems)
func (s *Sockets) sendMessageOut(msg models.Message) bool {
	keepGoing := true
	if msg.Type == models.COMMAND && msg.Body == "shutdown" {
		keepGoing = false
	}
	b, err := json.Marshal(msg)
	if err != nil {
		log.Println("zmq sockets had a problem", err)
	} else {
		s.pubChan.Out() <- [][]byte{
			[]byte(msg.Type),
			b,
		}
	}
	return keepGoing
}

//A message that came from outside clients (ui, connected
//gogadget systems) is passed along to this gogadget
//system
func (s *Sockets) sendMessageIn(data [][]byte, out chan<- models.Message) {
	if len(data) == 2 {
		msg := &models.Message{}
		json.Unmarshal(data[1], msg)
		msg.Sender = "zmq sockets"
		out <- *msg
	} else {
		log.Println("zmq received an improper message", data)
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
		err = s.getMasterSockets()
	} else {
		err = s.getClientSockets()
		s.subChan = s.sub.Channels()
		s.pubChan = s.pub.Channels()
	}
	return err
}

//Two GoGadgets systems can be joined together into a single system.
//The GoGadgets system that has it's App.Host set as "localhost" is
//the master.  Other systems that wish to join need to be configured
//with the IP address of the master system as App.Host.
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
//system or a UI.
func (s *Sockets) getClientSockets() (err error) {
	s.ctx, err = zmq.NewContext()
	if err != nil {
		return err
	}
	s.pub, err = s.ctx.Socket(zmq.Pub)
	if err != nil {
		return err
	}
	u := fmt.Sprintf("tcp://%s:%d", s.host, s.subPort)
	if err = s.pub.Connect(u); err != nil {
		return err
	}

	s.sub, err = s.ctx.Socket(zmq.Sub)
	if err != nil {
		return err
	}
	u = fmt.Sprintf("tcp://%s:%d", s.host, s.pubPort)
	if err = s.sub.Connect(u); err != nil {
		return err
	}
	s.sub.Subscribe([]byte(""))
	return err
}

func (s *Sockets) GetUID() string {
	return "zmq sockets"
}
