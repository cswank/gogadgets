package gogadgets

import (
	"labix.org/v2/mgo"
)

//Recorder takes all the update messages it receives and saves them
//in a mongodb.
type Recorder struct {
	OutputDevice
	DBHost string
	DBName string
	session *mgo.Session
	collection *mgo.Collection
	status bool
	connected bool
}

func NewRecorder(pin *Pin) (OutputDevice, error) {
	return &Recorder{
		DBHost: pin.Args["host"],
		DBName: pin.Args["db"],
	}, nil
}

func (r *Recorder) Update(msg *Message) {
	if msg.Type == "update" {
		r.save(msg)
	}
}

func (r *Recorder) On(val *Value) error {
	err := r.connect()
	r.status = true
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

func (r *Recorder)save(msg *Message) {
	r.collection.Insert(msg)
}

func (r *Recorder)connect() error {
	session, err := mgo.Dial(r.DBHost)
	if err != nil {
		return err
        }
	r.session = session
	r.collection = session.DB(r.DBName).C("updates")
	return nil
}
