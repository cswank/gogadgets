package gogadgets

import "strings"

type Cron struct {
	status bool
	jobs   jobs
}

func (c *Cron) Config() ConfigHelper {
	return ConfigHelper{}
}

func (c *Cron) Update(msg *Message) {

}

func (c *Cron) On(val *Value) error {
	c.status = true
	return nil
}

func (c *Cron) Status() interface{} {
	return c.status
}

func (c *Cron) Off() error {
	c.status = false
	return nil
}

type jobs struct {
	jobs map[string][]string
}

func (j *jobs) parse(s string) map[string]string {
	rows := strings.Split(s, "\n")
	for _, row := range rows {
		j.parseRow(row)
	}
}

func (j *jobs) parseRow(row string) {

}
