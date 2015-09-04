package gogadgets

import (
	"strings"
	"time"
)

type nower func() time.Time

func NewCron(pin *Pin) (Device, error) {
	return &Cron{
		jobs: pin.Args["jobs"].(string),
		Now:  time.Now,
	}, nil
}

type Cron struct {
	Now    nower
	status bool
	jobs   string
}

func (c *Cron) Start(in <-chan Message, out chan<- Message) {
	for {
		select {
		case t := <-time.After(time.Second):
			if t.Second() == 0 {
				c.checkJobs(t)
			}
		case msg := <-in:
			c.readMessage(msg)
		}
	}
}

func (c *Cron) checkJobs(t time.Time) {

}

//add new cron jobs via message?
func (c *Cron) readMessage(msg Message) {

}

type jobs struct {
	jobs map[string][]string
}

func (j *jobs) parse(s string) map[string]string {
	rows := strings.Split(s, "\n")
	for _, row := range rows {
		j.parseRow(row)
	}
	return nil
}

func (j *jobs) parseRow(row string) {

}
