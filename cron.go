package gogadgets

import (
	"fmt"
	"strings"
	"time"
)

type Afterer func(d time.Duration) <-chan time.Time

func NewCron(pin *Pin) (Device, error) {
	return &Cron{
		Jobs:  pin.Args["jobs"].(string),
		After: time.After,
		Sleep: time.Second,
	}, nil
}

type Cron struct {
	After  Afterer
	Jobs   string
	Sleep  time.Duration
	status bool
	jobs   map[string][]string
	out    chan<- Message
	ts     *time.Time
}

func (c *Cron) Start(in <-chan Message, out chan<- Message) {
	c.out = out
	c.parseJobs()
	for {
		select {
		case t := <-c.After(c.getSleep()):
			ts := time.Now()
			c.ts = &ts
			if t.Second() == 0 {
				c.checkJobs(t)
			}
		case msg := <-in:
			c.readMessage(msg)
		}
	}
}

func (c *Cron) getSleep() time.Duration {
	if c.ts == nil {
		return c.Sleep
	}
	diff := time.Now().Sub(*c.ts)
	return c.Sleep - diff
}

func (c *Cron) parseJobs() {
	c.jobs = map[string][]string{}
	rows := strings.Split(c.Jobs, "\n")
	for _, row := range rows {
		c.parseJob(row)
	}
}

func (c *Cron) parseJob(row string) {
	parts := strings.Split(row, " ")
	if len(parts) < 6 {
		return
	}
	key := c.getKey(parts[0:5])
	cmd := strings.Join(parts[5:], " ")
	a, ok := c.jobs[key]
	if !ok {
		a = []string{}
	}
	a = append(a, cmd)
	c.jobs[key] = a
}

func (c *Cron) getKey(parts []string) string {
	out := make([]string, len(parts))
	for i, part := range parts {
		out[i] = strings.Replace(part, " ", "", -1)
	}
	return strings.Join(parts, " ")
}

func (c *Cron) checkJobs(t time.Time) {
	keys := c.getKeys(t)
	for _, k := range keys {
		cmds, ok := c.jobs[k]
		if ok {
			for _, cmd := range cmds {
				c.out <- Message{
					Type:   COMMAND,
					Sender: "cron",
					UUID:   GetUUID(),
					Body:   cmd,
				}
			}
		}
	}
}

func (c *Cron) getKeys(t time.Time) []string {
	return []string{
		"* * * * *",
		fmt.Sprintf("%d * * * *", t.Minute()),
		fmt.Sprintf("%d %d * * *", t.Minute(), t.Hour()),
		fmt.Sprintf("%d %d %d * *", t.Minute(), t.Hour(), t.Day()),
		fmt.Sprintf("%d %d %d %d *", t.Minute(), t.Hour(), t.Day(), t.Month()),
		fmt.Sprintf("%d %d %d %d %d", t.Minute(), t.Hour(), t.Day(), t.Month(), t.Year()),
	}
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
