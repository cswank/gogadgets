package gogadgets

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/cswank/gogadgets/lib/cron"
)

var (
	cronExp, _ = regexp.Compile("^[*0-9,-]+$")
)

type Afterer func(d time.Duration) <-chan time.Time

func NewCron(config *GadgetConfig, options ...func(*Cron) error) (*Cron, error) {

	c := &Cron{
		client: &http.Client{Timeout: 5 * time.Second},
	}

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

	if config != nil && c.guards == nil {
		if raw, ok := config.Args["guards"]; ok {
			data, err := json.Marshal(raw)
			if err != nil {
				return c, fmt.Errorf("encode guards: %w", err)
			}
			if err := json.Unmarshal(data, &c.guards); err != nil {
				return c, fmt.Errorf("decode guards: %w", err)
			}
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

// CronGuards installs moisture-gate rules. Use from tests or callers that
// build a Cron programmatically. Config-driven setup populates the same
// field via "guards" in Args.
func CronGuards(g []Guard) func(*Cron) error {
	return func(c *Cron) error {
		c.guards = g
		return nil
	}
}

// CronClient overrides the HTTP client used to fetch guard readings.
// Intended for tests; production callers can rely on the default.
func CronClient(client *http.Client) func(*Cron) error {
	return func(c *Cron) error {
		c.client = client
		return nil
	}
}

// Guard skips a scheduled command when a remote sensor reading is at or
// above Max. Commands are matched by prefix so a Match of
// "turn on garden bed sprinklers" gates both the bare command and any
// "... for N minutes" variant.
type Guard struct {
	Match  string  `json:"match"`
	URL    string  `json:"url"`
	Device string  `json:"device"`
	Max    float64 `json:"max"`
}

type Cron struct {
	after  Afterer
	sleep  time.Duration
	status bool
	out    chan<- Message
	ts     *time.Time
	cron   *cron.Cron
	guards []Guard
	client *http.Client
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
	diff := time.Since(*c.ts)
	return c.sleep - diff
}

func (c *Cron) checkJobs(t time.Time) {
	cmds := c.cron.Check(t)
	for _, cmd := range cmds {
		if c.shouldSkip(cmd) {
			continue
		}
		c.out <- Message{
			Type:   COMMAND,
			Sender: "cron",
			UUID:   GetUUID(),
			Body:   cmd,
		}
	}
}

// shouldSkip returns true if any guard applies to cmd and reports a sensor
// value at or above its Max. Guard fetch errors fail open — irrigation
// should run when the sensor is unreachable rather than let plants die.
func (c *Cron) shouldSkip(cmd string) bool {
	for _, g := range c.guards {
		if !strings.HasPrefix(cmd, g.Match) {
			continue
		}
		v, err := c.fetchValue(g.URL, g.Device)
		if err != nil {
			log.Printf("cron: guard fetch for %q failed, proceeding: %s", g.Device, err)
			return false
		}
		if v >= g.Max {
			log.Printf("cron: skipping %q (%s reads %.1f, threshold %.1f)", cmd, g.Device, v, g.Max)
			return true
		}
	}
	return false
}

func (c *Cron) fetchValue(url, device string) (float64, error) {
	resp, err := c.client.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("status %d", resp.StatusCode)
	}
	var msgs []Message
	if err := json.NewDecoder(resp.Body).Decode(&msgs); err != nil {
		return 0, err
	}
	for _, m := range msgs {
		if m.Location+" "+m.Name != device {
			continue
		}
		f, ok := m.Value.ToFloat()
		if !ok {
			return 0, fmt.Errorf("non-numeric value for %q", device)
		}
		return f, nil
	}
	return 0, fmt.Errorf("device %q not found", device)
}
