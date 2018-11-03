package gogadgets

import (
	"regexp"
	"time"

	"github.com/cswank/gogadgets/lib/cron"
)

var (
	cronExp, _ = regexp.Compile("^[*0-9,-]+$")
)

type Afterer func(d time.Duration) <-chan time.Time

func NewCron(config *GadgetConfig, options ...func(*Cron) error) (*Cron, error) {

	c := &Cron{}

	for _, opt := range options {
		if err := opt(c); err != nil {
			return c, err
		}
	}

	if c.after == nil {
		c.after = time.After
	}

	if c.sleep == time.Duration(0) {
		c.sleep = time.Second
	}

	var err error
	if c.cron == nil {
		v := config.Args["jobs"].([]interface{})
		jobs := make([]string, len(v))
		for i, r := range v {
			jobs[i] = r.(string)
		}
		c.cron, err = cron.New(jobs)
		if err != nil {
			return c, err
		}
	}

	return c, nil
}

func CronAfter(a Afterer) func(*Cron) error {
	return func(c *Cron) error {
		c.after = a
		return nil
	}
}

func CronSleep(d time.Duration) func(*Cron) error {
	return func(c *Cron) error {
		c.sleep = d
		return nil
	}
}

func CronJobs(j []string) func(*Cron) error {
	return func(c *Cron) error {
		var err error
		c.cron, err = cron.New(j)
		return err
	}
}

type Cron struct {
	after  Afterer
	sleep  time.Duration
	status bool
	out    chan<- Message
	ts     *time.Time
	cron   *cron.Cron
}

func (c *Cron) GetUID() string {
	return "cron"
}

func (c *Cron) GetDirection() string {
	return "na"
}

func (c *Cron) Start(in <-chan Message, out chan<- Message) {
	c.out = out
	for {
		select {
		case t := <-c.after(c.getSleep()):
			ts := time.Now()
			c.ts = &ts
			if t.Second() == 0 {
				c.checkJobs(t)
			}
		case <-in:
		}
	}
}

func (c *Cron) getSleep() time.Duration {
	if c.ts == nil {
		return c.sleep
	}
	diff := time.Now().Sub(*c.ts)
	return c.sleep - diff
}

func (c *Cron) checkJobs(t time.Time) {
	cmds := c.cron.Check(t)
	for _, cmd := range cmds {
		c.out <- Message{
			Type:   COMMAND,
			Sender: "cron",
			UUID:   GetUUID(),
			Body:   cmd,
		}
	}
}
