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
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/vaughan0/go-zmq"
)

//Sockets fufills the GoGadget interface and is
//added to each Gadget system by App.  It provides
//a way to connect multiple gadgets systems together
//as a single system, and also provides a way for
//an external UI to interact the system.
type Sockets struct {
	id      string
	master  bool
	host    string
	pubPort int
	subPort int
	ctx     *zmq.Context
	sub     *zmq.Socket
	subChan *zmq.Channels
	pub     *zmq.Socket
	pubChan *zmq.Channels
	updates map[string]Message
}

type SocketsConfig struct {
	Host    string
	SubPort int
	PubPort int
	Master  bool
}

func NewSockets(cfg SocketsConfig) *Sockets {
	if cfg.Master {
		return newMasterSockets(cfg)
	}
	return NewClientSockets(cfg)
}

func newMasterSockets(cfg SocketsConfig) *Sockets {
	return &Sockets{
		id:      newUUID(),
		master:  true,
		host:    cfg.Host,
		subPort: cfg.SubPort,
		pubPort: cfg.PubPort,
		updates: map[string]Message{},
	}
}

func NewClientSockets(cfg SocketsConfig) *Sockets {
	return &Sockets{
		master:  false,
		id:      newUUID(),
		host:    cfg.Host,
		subPort: cfg.SubPort,
		pubPort: cfg.PubPort,
	}
}

func (s *Sockets) Send(cmd string) {
	msg := Message{
		From: s.id,
		Type: COMMAND,
		Body: cmd,
	}
	s.SendMessage(msg)
}

func (s *Sockets) SendMessage(msg Message) {
	if s.pubChan == nil {
		s.Connect()
	}
	msg.From = s.id
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

func (s *Sockets) Recv() *Message {
	if s.subChan == nil {
		s.Connect()
	}
	data, err := s.sub.Recv()
	if err != nil {
		panic(err)
	}
	msg := &Message{}
	json.Unmarshal(data[1], msg)
	if msg.From == s.id {
		return nil
	}
	return msg
}

func (s *Sockets) SendStatusRequest() (map[string]Message, error) {
	msgs := map[string]Message{}
	tries := 0
	s.SendMessage(Message{
		Body: "status",
		From: s.id,
	})
	for {
		data, err := s.sub.Recv()
		if err != nil {
			panic(err)
		}
		if string(data[0]) == "status" {
			err = json.Unmarshal(data[1], &msgs)
			return msgs, err
		} else {
			tries += 1
		}
		if tries > 5 {
			return msgs, errors.New("didn't get a status response")
		}
	}
	return msgs, nil
}

//Sockets listens for chann Messages from inside the system and
//sends it to external listeners (like a UI), and listens for
//external messages and sends them along to the internal system.
func (s *Sockets) Start(in <-chan Message, out chan<- Message) {
	if err := s.Connect(); err != nil {
		lg.Fatal(err)
	}
	for {
		select {
		case data := <-s.subChan.In():
			msg, err := s.unmarshalMessage(data)
			if err == nil && msg.From != s.id {
				s.sendMessageIn(msg, out)
			}
			if s.master && msg.From == s.id {
				s.sendMessageOut(*msg)
			}
		case msg := <-in:
			if msg.From == "" {
				msg.From = s.id
				s.sendMessageOut(msg)
			}
		case err := <-s.subChan.Errors():
			log.Println(err)
		}
	}
}

//A message that came from inside this gogadgets system
//is sent to outside clients (ui, connected gogadget systems)
func (s *Sockets) sendMessageOut(msg Message) bool {
	if s.master && msg.Type == UPDATE {
		s.updates[msg.Sender] = msg
	}
	keepGoing := true
	if msg.Type == COMMAND && msg.Body == "shutdown" {
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

func (s *Sockets) unmarshalMessage(data [][]byte) (*Message, error) {
	var msg *Message
	if len(data) != 2 {
		return msg, errors.New("improper message")
	}
	msg = &Message{}
	return msg, json.Unmarshal(data[1], msg)
}

//A message that came from outside clients (ui, connected
//gogadget systems) is passed along to this gogadget
//system
func (s *Sockets) sendMessageIn(msg *Message, out chan<- Message) {
	if msg.Sender == "" {
		msg.Sender = "zmq sockets"
	}
	if s.master && msg.Type == UPDATE {
		s.updates[msg.Sender] = *msg
	}
	if msg.Body == "status" {
		s.sendStatus()
	} else {
		out <- *msg
	}
}

//An outside client (like a UI) wants the latest status of
//all gadgets in the system.
func (s *Sockets) sendStatus() {
	b, _ := json.Marshal(&s.updates)
	s.pubChan.Out() <- [][]byte{
		[]byte("status"),
		b,
	}
}

func (s *Sockets) Close() {
	s.sub.Close()
	s.pub.Close()
	if s.subChan != nil {
		s.subChan.Close()
	}
	if s.pubChan != nil {
		s.pubChan.Close()
	}
	s.ctx.Close()
}

func (s *Sockets) Connect() (err error) {
	if s.master {
		err = s.getMasterChannels()
	} else {
		err = s.getClientSockets()
		s.subChan = s.sub.Channels()
		s.pubChan = s.pub.Channels()
	}
	time.Sleep(100 * time.Millisecond)
	return err
}

//Two GoGadgets systems can be joined together into a single system.
//The GoGadgets system that has it's App.Host set as "localhost" is
//the master.  Other systems that wish to join need to be configured
//with the IP address of the master system as App.Host.
func (s *Sockets) getMasterChannels() (err error) {
	s.master = true
	err = s.getMasterSockets()
	if err != nil {
		return err
	}
	s.pubChan = s.pub.Channels()
	s.subChan = s.sub.Channels()
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
	sub, err := s.ctx.Socket(zmq.Sub)
	if err != nil {
		return err
	}
	s.sub = sub
	if err = s.sub.Bind(fmt.Sprintf("tcp://*:%d", s.subPort)); err != nil {
		return err
	}
	s.sub.Subscribe([]byte(""))
	return err
}

//This creates the zmq sockets for a GoGadget system that is not the master
//system or a UI.
func (s *Sockets) getClientSockets() (err error) {
	s.master = false
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
