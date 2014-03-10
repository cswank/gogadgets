package main

import (
	"flag"
	"labix.org/v2/mgo"
	"bitbucket.org/cswank/gogadgets"
)

var (
	cfg = flag.String("c", "", "Path to the config json file")
	db = flag.String("d", "gadgets", "database name")
)

//Recorder takes all the update messages it receives and saves them
//in a mongodb.
type Recorder struct {
	gogadgets.GoGadget
	session *mgo.Session
	collection *mgo.Collection
}

func (r *Recorder)Start(in <-chan gogadgets.Message, out chan<- gogadgets.Message) {
	r.connect()
	defer r.session.Close()
	for {
		msg := <-in
		if msg.Type == "update" {
			r.save(&msg)
		}
	}
}

func (r *Recorder)save(msg *gogadgets.Message) {
	r.collection.Insert(msg)
}

func (r *Recorder)GetUID() string {
	return "recorder"
}

func main() {
	flag.Parse()
	a := gogadgets.NewApp(*cfg)
	r := &Recorder{}
	a.AddGadget(r)
	a.Start()
}

func (r *Recorder)connect() error {
	session, err := mgo.Dial("localhost")
	if err != nil {
		return err
        }
	r.session = session
	r.collection = session.DB(*db).C("updates")
	return nil
}
