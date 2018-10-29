package gogadgets

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type SMS struct {
	//twillio sid and oauth
	url     string
	sid     string
	token   string
	from    string
	message string
	to      []string
}

func NewSMS(pin *Pin) (OutputDevice, error) {
	sid, ok := pin.Args["sms"].(string)
	if !ok {
		return nil, fmt.Errorf("could not parse sms from pin args")
	}
	from, ok := pin.Args["from"].(string)
	if !ok {
		return nil, fmt.Errorf("could not parse from from pin args")
	}
	msg, ok := pin.Args["message"].(string)
	if !ok {
		return nil, fmt.Errorf("could not parse message from pin args")
	}
	token, ok := pin.Args["token"].(string)
	if !ok {
		return nil, fmt.Errorf("could not parse token from pin args")
	}
	var to []string
	tos, ok := pin.Args["to"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("could not parse to from pin args")
	}
	for _, v := range tos {
		val, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("could not parse to from pin args")
		}
		tos = append(tos, val)
	}

	return &SMS{
		sid:     sid,
		from:    from,
		message: msg,
		url:     fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", sid),
		token:   token,
		to:      to,
	}, nil
}

func (s *SMS) WithOutput(out map[string]OutputDevice) {}

func (s *SMS) Commands(location, name string) *Commands {
	return nil
}

func (s *SMS) Config() ConfigHelper {
	return ConfigHelper{}
}

func (s *SMS) Update(msg *Message) bool {
	return false
}

func (s *SMS) On(val *Value) error {
	msgData := url.Values{}
	for _, to := range s.to {
		msgData.Add("To", to)
	}
	msgData.Set("From", s.from)
	msgData.Set("Body", s.message)
	msgDataReader := *strings.NewReader(msgData.Encode())
	client := &http.Client{}
	req, _ := http.NewRequest("POST", s.url, &msgDataReader)
	req.SetBasicAuth(s.sid, s.token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, _ := client.Do(req)
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			log.Println(data["sid"])
		}
	} else {
		log.Println(resp.Status)
	}
	return nil
}

func (s *SMS) Status() map[string]bool {
	return map[string]bool{}
}

func (s *SMS) Off() error {
	return nil
}
