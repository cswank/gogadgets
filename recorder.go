package gogadgets

import (
	"encoding/json"
	"labix.org/v2/mgo"
	"log"
	"time"
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
	retries    int
}

func NewRecorder(pin *Pin) (OutputDevice, error) {
	s := getSummaries(pin.Args["summarize"])
	r := &Recorder{
		DBHost:    getDbHost(pin.Args["host"]),
		DBName:    getDbName(pin.Args["db"]),
		filter:    getFilter(pin.Args["filter"]),
		history:   map[string]summary{},
		summaries: s,
	}
	return r, nil
}

func getDbHost(host interface{}) string {
	if host == nil {
		return "localhost"
	}
	return host.(string)
}

func getDbName(db interface{}) string {
	if db == nil {
		return "gogadgets"
	}
	return db.(string)
}

func (r *Recorder) Config() ConfigHelper {
	return ConfigHelper{
		Args: map[string]interface{}{
			"host": []string{},
			"db":   []string{},
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
	err := r.connect()
	if err == nil {
		r.status = true
	}
	return err
}

func (r *Recorder) Off() error {
	if r.session != nil {
		r.close()
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
	err := r.write(msg)
	if err != nil {
		if r.retries > 5 {
			log.Println("couldn't connect to the db")
			return
		}
		r.close()
		err = r.connect()
		if err != nil {
			log.Println("couldn't connect to the db")
			return
		}
		r.doSave(msg)
		r.retries += 1
	} else {
		r.retries = 0
	}
}

func (r *Recorder) close() {
	if r.session != nil {
		r.session.Close()
	}
	r.session = nil
	r.collection = nil
}

func (r *Recorder) write(msg *Message) error {
	var err error
	if r.collection == nil {
		if err = r.connect(); err != nil {
			return err
		}
	}
	return r.collection.Insert(msg)
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
	if f == nil {
		return []string{}
	}
	filters, ok := f.([]string)
	if !ok {
		return []string{}
	}
	return filters
}
