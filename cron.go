package gogadgets

import (
	"bytes"
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"
)

type Afterer func(d time.Duration) <-chan time.Time

func NewCron(config *GadgetConfig) (*Cron, error) {
	v := config.Args["jobs"].([]interface{})
	jobs := make([]string, len(v))
	for i, r := range v {
		jobs[i] = r.(string)
	}
	return &Cron{
		Jobs:  jobs,
		After: time.After,
		Sleep: time.Second,
	}, nil
}

type Cron struct {
	After  Afterer
	Jobs   []string
	Sleep  time.Duration
	status bool
	jobs   map[string][]string
	out    chan<- Message
	ts     *time.Time
}

func (c *Cron) GetUID() string {
	return "cron"
}

func (c *Cron) GetDirection() string {
	return "na"
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
		case <-in:
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
	for _, row := range c.Jobs {
		c.parseJob(row)
	}
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
		lg.Printf("could not parse %s", s)
		return []string{}
	}
	start, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		lg.Printf("could not parse %s", s)
		return []string{}
	}
	end, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		lg.Printf("could not parse %s", s)
		return []string{}
	}
	if end <= start {
		lg.Printf("could not parse %s", s)
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

type now struct {
	Minute  int
	Hour    int
	Day     int
	Month   int
	Weekday int
}

func (c *Cron) getPossibilities(t time.Time) []string {
	n := now{
		Minute:  t.Minute(),
		Hour:    t.Hour(),
		Day:     t.Day(),
		Month:   int(t.Month()),
		Weekday: int(t.Weekday()),
	}
	tpl, _ := template.New("possibilites").Parse(`
* * * * *
{{.Minute}} * * * *
* {{.Hour}} * * *
* * {{.Day}} * *
* * * {{.Month}} *
* * * * {{.Weekday}}
{{.Minute}} {{.Hour}} * * *
{{.Minute}} * {{.Day}} * *
{{.Minute}} * * {{.Month}} *
{{.Minute}} * * * {{.Weekday}}
* {{.Hour}} {{.Day}} * *
* * {{.Day}} {{.Month}} *
* * * {{.Month}} {{.Weekday}}
* {{.Hour}} * {{.Month}} *
* {{.Hour}} * * {{.Weekday}}
* * {{.Day}} * {{.Weekday}}
* * {{.Day}} {{.Month}} {{.Weekday}}
* {{.Hour}} * {{.Month}} {{.Weekday}}
* {{.Hour}} {{.Day}} * {{.Weekday}}
* {{.Hour}} {{.Day}} {{.Month}} *
{{.Minute}} * * {{.Month}} {{.Weekday}}
{{.Minute}} {{.Hour}} * * {{.Weekday}}
{{.Minute}} {{.Hour}} {{.Day}} * *
{{.Minute}} * {{.Day}} * {{.Weekday}}
{{.Minute}} * {{.Day}} {{.Month}} *
{{.Minute}} {{.Hour}} * {{.Month}} *
{{.Minute}} {{.Hour}} {{.Day}} {{.Month}} *
{{.Minute}} {{.Hour}} {{.Day}} * {{.Weekday}}
{{.Minute}} {{.Hour}} * {{.Month}} {{.Weekday}}
{{.Minute}} * {{.Day}} {{.Month}} {{.Weekday}}
* {{.Minute}} {{.Hour}} {{.Day}} {{.Weekday}}
{{.Minute}} {{.Hour}} {{.Day}} {{.Month}} {{.Weekday}}
`)
	buf := bytes.Buffer{}
	tpl.Execute(&buf, n)
	return strings.Split(buf.String(), "\n")
}
