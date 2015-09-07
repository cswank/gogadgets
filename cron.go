package gogadgets

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

type Afterer func(d time.Duration) <-chan time.Time

func NewCron(config *GadgetConfig) (*Cron, error) {
	return &Cron{
		Jobs:  config.Args["jobs"].(string),
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

func (c *Cron) GetUID() string {
	return "cron"
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

func (c *Cron) readMessage(msg Message) {

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
	fmt.Println(c.jobs)
}

func (c *Cron) parseJob(row string) {
	if strings.Index(row, "#") == 0 {
		return
	}
	parts := strings.Fields(row)
	if len(parts) < 6 {
		return
	}
	keys := c.getKeys(parts[0:5])
	cmd := strings.Join(parts[5:], " ")
	for _, key := range keys {
		a, ok := c.jobs[key]
		if !ok {
			a = []string{}
		}
		a = append(a, cmd)
		c.jobs[key] = a
	}
}

func (c *Cron) getKeys(parts []string) []string {
	out := []string{}
	var hasRange bool
	for i, part := range parts {
		if strings.Index(part, "-") >= 1 {
			hasRange = true
			r := c.getRange(part)
			for _, x := range r {
				parts[i] = x
				out = append(out, c.getKeys(parts)...)
			}
		} else if strings.Index(part, ",") >= 1 {
			hasRange = true
			s := strings.Split(part, ",")
			for _, x := range s {
				parts[i] = x
				out = append(out, c.getKeys(parts)...)
			}
		}
	}
	if !hasRange {
		out = append(out, strings.Join(parts, " "))
	}
	return out
}

func (c *Cron) getRange(s string) []string {
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		lg.Printf("could not parse %", s)
		return []string{}
	}
	start, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		lg.Printf("could not parse %", s)
		return []string{}
	}
	end, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		lg.Printf("could not parse %", s)
		return []string{}
	}
	if end <= start {
		lg.Printf("could not parse %", s)
		return []string{}
	}
	out := make([]string, end-start+1)
	j := 0
	for i := start; i <= end; i++ {
		out[j] = fmt.Sprintf("%d", i)
		j++
	}
	return out
}

func (c *Cron) checkJobs(t time.Time) {
	keys := c.getPossibilities(t)
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

func (c *Cron) getPossibilities(t time.Time) []string {
	parts := [][]string{
		[]string{"*", fmt.Sprintf("%d", t.Minute())},
		[]string{"*", fmt.Sprintf("%d", t.Hour())},

		[]string{"*", fmt.Sprintf("%d", t.Day())},
		[]string{"*", fmt.Sprintf("%d", t.Month())},
		[]string{"*", fmt.Sprintf("%d", t.Weekday())},
	}
	f := `
* * * * *
%s * * * *
* %s * * *
* * %s * *
* * * %s *
* * * * %s
%s %s * * *
%s * %s * *
%s * * %s *
%s * * * %s
* %s %s * *
* * %s %s *
* * * %s %s
* %s * %s *
* %s * * %s
* * %s * %s
* * %s %s %s
* %s * %s %s
* %s %s * %s
* %s %s %s *
%s * * %s %s
%s %s * * %s
%s %s %s * *
%s * %s * %s
%s * %s %s *
%s %s * %s *
%s %s %s %s *
%s %s %s * %s
%s %s * %s %s
%s * %s %s %s
* %s %s %s %s
%s %s %s %s %s
`

	//"* * * * *"
	a := make([][]string, 25)
	for i := range a {
		a[i] = make([]string, 5)
	}
	for r := 0; r < 25; r++ {
		for c := 0; c < 5; c++ {
			a[r][c] = parts[c][((r + c) % 2)]
		}
	}
	log.Fatal(a)
	return []string{}
}
