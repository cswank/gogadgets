package gogadgets

import (
	"labix.org/v2/mgo"
	"encoding/json"
	"time"
	"log"
)

type summary struct {
	start time.Time
	n     int
	v     float64
}

//Recorder takes all the update messages it receives and saves them
//in a mongodb.
type Recorder struct {
	DBHost     string
	DBName     string
	session    *mgo.Session
	collection *mgo.Collection
	status     bool
	connected  bool
	filter     []string
	summaries  map[string]time.Duration
	history    map[string]summary
}

func NewRecorder(pin *Pin) (OutputDevice, error) {
	r := &Recorder{
		DBHost:    pin.Args["host"].(string),
		DBName:    pin.Args["db"].(string),
		filter:    getFilter(pin.Args["filter"]),
		history:   map[string]summary{},
		summaries: getSummaries(pin.Args["summarize"]),
	}
	return r, nil
}

func getSummaries(s interface{}) map[string]time.Duration {
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
	err := r.connect()
	if err == nil {
		r.status = true
	}
	return err
}

func (r *Recorder) Off() error {
	if r.session != nil {
		r.session.Close()
	}
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
	r.collection.Insert(msg)
}

func (r *Recorder) connect() error {
	session, err := mgo.Dial(r.DBHost)
	if err != nil {
		return err
	}
	r.session = session
	r.collection = session.DB(r.DBName).C("updates")
	return nil
}

func getFilter(f interface{}) []string {
	filters, ok := f.([]string)
	if !ok {
		return []string{}
	}
	return filters
}
