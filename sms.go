package gogadgets

import (
	"log"

	"github.com/subosito/twilio"
)

// Sends text messages (der)
type SMS struct {
	status       bool
	trigger      string
	triggerState bool
	Driver       Driver
}

func NewSMS(pin *Pin) (OutputDevice, error) {
	args := pin.Args
	t := args["trigger"].(string)
	s := args["triggerState"].(bool)
	var d Driver
	if args["driver"] == "twilio" {
		d = NewTwilioDriver(args)
	} else {
		log.Fatal("could not create an sms driver (twilio)")
	}
	return &SMS{
		trigger:      t,
		triggerState: s,
		Driver:       d,
	}, nil
}

func (s *SMS) Config() ConfigHelper {
	return ConfigHelper{}
}

func (s *SMS) Update(msg *Message) {
	if s.shouldSend(msg) {
		if err := s.Driver.Send(); err != nil {
			log.Println("unable to send sms", err)
		}
	}
}

func (s *SMS) shouldSend(msg *Message) bool {
	return s.status &&
		msg.Type == "update" &&
		msg.Sender == s.trigger &&
		msg.Value.Value.(bool) == s.triggerState
}

func (s *SMS) On(val *Value) error {
	s.status = true
	return nil
}

func (s *SMS) Status() interface{} {
	return s.status
}

func (s *SMS) Off() error {
	s.status = false
	return nil
}

type Driver interface {
	Send() error
}

type Twilio struct {
	msg    string
	to     string
	from   string
	cli    *twilio.Client
	params twilio.MessageParams
}

func NewTwilioDriver(args map[string]interface{}) *Twilio {
	sid := args["account_sid"].(string)
	token := args["auth_token"].(string)
	from := args["from"].(string)
	to := args["to"].(string)
	body := args["body"].(string)
	return &Twilio{
		from: from,
		cli:  twilio.NewClient(sid, token, nil),
		to:   to,
		params: twilio.MessageParams{
			Body: body,
		},
	}
}

func (t *Twilio) Send() error {
	_, _, err := t.cli.Messages.Send(t.from, t.to, t.params)
	return err
}
