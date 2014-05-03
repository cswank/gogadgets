package gogadgets

import (
	"bitbucket.org/cswank/gogadgets/models"
)

//All the gadgets of the system push their messages here.
type Broker struct {
	queue    *Queue
	channels map[string]chan models.Message
	collect  chan models.Message
	input    <-chan models.Message
}

func NewBroker(channels map[string]chan models.Message, input <-chan models.Message, collect chan models.Message) *Broker {
	return &Broker{
		input:    input,
		queue:    NewQueue(),
		channels: channels,
		collect:  collect,
	}
}

func (b *Broker) Start() {
	in := make(chan models.Message)
	go b.collectMessages(b.collect)
	go b.dispenseMessages(in)
	keepRunning := true
	for keepRunning {
		select {
		case msg := <-in:
			b.sendMessage(msg)
		case msg := <-b.input:
			if msg.Type == "command" && msg.Body == "shutdown" {
				keepRunning = false
			}
			b.sendMessage(msg)
		}
	}
}

//Collects each message that is sent by the parts of the
//system and pushes it in the queue.
func (b *Broker) collectMessages(in <-chan models.Message) {
	for {
		msg := <-in
		b.queue.Push(&msg)
	}
}

//After a message is collected by collectMessage, it is
//then sent back to the rest of the system.
func (b *Broker) dispenseMessages(out chan<- models.Message) {
	for {
		b.queue.Lock()
		if b.queue.Len() == 0 {
			b.queue.Wait()
		}
		msg := b.queue.Get()
		out <- *msg
		b.queue.Unlock()
	}
}

func (b *Broker) sendMessage(msg models.Message) {
	if msg.Target == "" {
		for uid, channel := range b.channels {
			if uid != msg.Sender {
				channel <- msg
			}
		}
	} else {
		channel, ok := b.channels[msg.Target]
		if ok {
			channel <- msg
		}
	}
}
