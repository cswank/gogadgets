package gogadgets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type summary struct {
	start time.Time
	n     int
	v     float64
}

//Recorder takes all the update messages it receives and saves them
//by posting to quimby
type Recorder struct {
	url       string
	token     string
	status    bool
	filter    []string
	summaries map[string]time.Duration
	history   map[string]summary
}

func NewRecorder(pin *Pin) (OutputDevice, error) {
	s := getSummaries(pin.Args["summarize"])
	r := &Recorder{
		url:       pin.Args["host"].(string),
		token:     pin.Args["token"].(string),
		filter:    getFilter(pin.Args["filter"]),
		history:   map[string]summary{},
		summaries: s,
	}
	return r, nil
}

func (r *Recorder) Config() ConfigHelper {
	return ConfigHelper{
		Args: map[string]interface{}{
			"host": []string{},
		},
	}
}

func getSummaries(s interface{}) map[string]time.Duration {
	if s == nil {
		return map[string]time.Duration{}
	}
	d, _ := json.Marshal(s)
	vals := map[string]int{}
	err := json.Unmarshal(d, &vals)
	out := map[string]time.Duration{}
	if err != nil {
		log.Println("WARNING, could not parse recorder summaires", s)
		return out
	}
	for key, val := range vals {
		var d time.Duration
		d = time.Duration(val) * time.Minute
		out[key] = d
	}
	return out
}

func (r *Recorder) Update(msg *Message) {
	if r.status && msg.Type == "update" {
		r.save(msg)
	}
}

func (r *Recorder) On(val *Value) error {
	r.status = true
	return nil
}

func (r *Recorder) Off() error {
	r.status = false
	return nil
}

func (r *Recorder) Status() interface{} {
	return r.status
}

func (r *Recorder) save(msg *Message) {
	if len(r.filter) > 0 {
		if !r.inFilter(msg) {
			return
		}
	}
	d, ok := r.summaries[msg.Sender]
	if ok {
		r.summarize(msg, d)
	} else {
		r.doSave(msg)
	}
}

func (r *Recorder) inFilter(msg *Message) bool {
	for _, item := range r.filter {
		if msg.Sender == item {
			return true
		}
	}
	return false
}

func (r *Recorder) summarize(msg *Message, duration time.Duration) {
	now := time.Now().UTC()
	s, ok := r.history[msg.Sender]
	if !ok {
		s = summary{start: now}
	}
	s.n += 1
	f, _ := msg.Value.ToFloat()
	s.v += f
	lapsed := now.Sub(s.start)
	if lapsed >= duration {
		msg.Value.Value = s.v / float64(s.n)
		r.doSave(msg)
		delete(r.history, msg.Sender)
	} else {
		r.history[msg.Sender] = s
	}
}

func (r *Recorder) doSave(msg *Message) {
	m := map[string]float64{
		"value": msg.Value.Value.(float64),
	}
	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	enc.Encode(m)

	u := fmt.Sprintf(r.url, msg.Location, msg.Name)
	req, err := http.NewRequest("POST", u, &buf)
	if err != nil {
		log.Println("couldn't post data", err)
		return
	}
	req.Header.Add("Authorization", r.token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("couldn't post data", err)
		return
	}
	resp.Body.Close()
}

func getFilter(f interface{}) []string {
	if f == nil {
		return []string{}
	}
	filters, ok := f.([]string)
	if !ok {
		return []string{}
	}
	return filters
}
