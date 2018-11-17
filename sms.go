package gogadgets

import (
	"fmt"
	"io/ioutil"
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
	sid, ok := pin.Args["sid"].(string)
	if !ok {
		return nil, fmt.Errorf("could not parse sid from pin args")
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
		to = append(to, val)
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

func (s *SMS) Commands(location, name string) *Commands {
	return nil
}

func (s *SMS) Update(msg *Message) bool {
	return false
}

func (s *SMS) On(val *Value) error {
	for _, to := range s.to {
		if err := s.sms(to); err != nil {
			return err
		}
	}

	return nil
}

func (s *SMS) sms(to string) error {
	req, err := s.request(to)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		d, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("error from twillio: %s", string(d))
	}

	return nil
}

func (s *SMS) request(to string) (*http.Request, error) {
	body := url.Values{
		"To":   []string{to},
		"From": []string{s.from},
		"Body": []string{s.message},
	}

	req, err := http.NewRequest("POST", s.url, strings.NewReader(body.Encode()))
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(s.sid, s.token)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

func (s *SMS) Status() map[string]bool {
	return map[string]bool{}
}

func (s *SMS) Off() error {
	return nil
}
