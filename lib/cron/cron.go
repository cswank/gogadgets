package cron

import (
	"bytes"
	"fmt"
	"html/template"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	cronExp, _ = regexp.Compile("^[*0-9,-]+$")
)

type Cron struct {
	jobs map[string][]string
}

func New(in []string) (*Cron, error) {
	jobs, err := parseJobs(in)
	return &Cron{
		jobs: jobs,
	}, err
}

func (c *Cron) Check(t time.Time) []string {
	keys := c.getPossibilities(t)
	var out []string
	for _, k := range keys {
		cmds := c.jobs[k]
		out = append(out, cmds...)
	}
	return out
}

func parseJobs(jobs []string) (map[string][]string, error) {
	m := map[string][]string{}
	for _, row := range jobs {
		if err := parseJob(row, m); err != nil {
			return m, err
		}
	}
	return m, nil
}

func parseJob(row string, m map[string][]string) error {
	if strings.Index(row, "#") == 0 {
		return nil
	}

	parts := strings.Fields(row)
	if len(parts) < 6 {
		return fmt.Errorf("could not parse job: %s", row)
	}

	for _, p := range parts[0:5] {
		if !cronExp.MatchString(p) {
			return fmt.Errorf("could not parse job: %s", row)
		}
	}

	keys, err := getKeys(parts[0:5])
	if err != nil {
		return err
	}

	cmd := strings.Join(parts[5:], " ")
	for _, key := range keys {
		a, ok := m[key]
		if !ok {
			a = []string{}
		}
		a = append(a, cmd)
		m[key] = a
	}
	return nil
}

func getKeys(parts []string) ([]string, error) {
	out := []string{}
	var hasRange bool
	for i, part := range parts {
		if strings.Index(part, "-") >= 1 {
			hasRange = true
			r, err := getRange(part)
			if err != nil {
				return nil, err
			}
			for _, x := range r {
				parts[i] = x
				keys, err := getKeys(parts)
				if err != nil {
					return nil, err
				}
				out = append(out, keys...)
			}
		} else if strings.Index(part, ",") >= 1 {
			hasRange = true
			s := strings.Split(part, ",")
			for _, x := range s {
				parts[i] = x
				keys, err := getKeys(parts)
				if err != nil {
					return nil, err
				}
				out = append(out, keys...)
			}
		}
	}
	if !hasRange {
		out = append(out, strings.Join(parts, " "))
	}
	return out, nil
}

func getRange(s string) ([]string, error) {
	parts := strings.Split(s, "-")
	if len(parts) != 2 {
		return nil, fmt.Errorf("could not parse cron string %s", s)
	}
	start, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return nil, err
	}
	end, err := strconv.ParseInt(parts[1], 10, 32)
	if err != nil {
		return nil, err
	}
	if end <= start {
		return nil, fmt.Errorf("could not parse %s (end is before start)", s)
	}

	out := make([]string, end-start+1)
	j := 0
	for i := start; i <= end; i++ {
		out[j] = fmt.Sprintf("%d", i)
		j++
	}

	return out, nil
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
