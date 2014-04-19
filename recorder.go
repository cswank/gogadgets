package gogadgets

import (
	"labix.org/v2/mgo"
	"strconv"
	"strings"
	"time"
)

type summary struct {
	start time.Time
	n int
	v float64
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
	duration   time.Duration
	history    map[string]summary
}

func NewRecorder(pin *Pin) (OutputDevice, error) {
	i, err := strconv.ParseInt(pin.Args["summarize"], 10, 64)
	var d time.Duration
	if err == nil {
		d = time.Duration(i) * time.Minute
	}
	r := &Recorder{
		DBHost: pin.Args["host"],
		DBName: pin.Args["db"],
		filter: getFilter(pin.Args["filter"]),
		duration: d,
		history: map[string]summary{},
	}
	return r, nil
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
	if r.duration > 0 {
		r.summarize(msg)
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

func (r *Recorder) summarize(msg *Message) {
	now := time.Now().UTC()
	s, ok := r.history[msg.Sender]
	if !ok {
		s = summary{start: now}
	}
	s.n += 1
	f, _ := msg.Value.ToFloat()
	s.v += f
	lapsed := s.start.Sub(now)
	if lapsed >= r.duration {
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

func getFilter(filterStr string) []string {
	if len(filterStr) == 0 {
		return []string{}
	}
	return strings.Split(filterStr, ",")
}
