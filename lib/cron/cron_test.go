package cron_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/cswank/gogadgets/lib/cron"
	. "github.com/onsi/gomega"
)

/*
January 2018
Su Mo Tu We Th Fr Sa
    1  2  3  4  5  6
 7  8  9 10 11 12 13
14 15 16 17 18 19 20
21 22 23 24 25 26 27
28 29 30 31
*/

func TestCron(t *testing.T) {
	g := NewGomegaWithT(t)

	testCases := []struct {
		doc  string
		cmds []string
		time string
		want []string
	}{
		//matches
		{"everything", []string{"* * * * * turn on living room light"}, "05 Jan 18 13:25 UTC", []string{"turn on living room light"}},
		{"minute", []string{"25 * * * * turn on living room light"}, "01 Jan 18 13:25 UTC", []string{"turn on living room light"}},
		{"minute and hour", []string{"25 13 * * * turn on living room light"}, "01 Jan 18 13:25 UTC", []string{"turn on living room light"}},
		{"minute, hour, and day", []string{"25 13 1 * * turn on living room light"}, "01 Jan 18 13:25 UTC", []string{"turn on living room light"}},
		{"minute, hour, day, month", []string{"25 13 1 1 * turn on living room light"}, "01 Jan 18 13:25 UTC", []string{"turn on living room light"}},
		{"weekday", []string{"25 13 * * 5 turn on living room light"}, "05 Jan 18 13:25 UTC", []string{"turn on living room light"}},
		{"multiple weekdays", []string{"25 13 * * 1,3,6 turn on living room light"}, "01 Jan 18 13:25 UTC", []string{"turn on living room light"}},
		{"range of weekdays", []string{"25 13 * * 4-6 turn on living room light"}, "04 Jan 18 13:25 UTC", []string{"turn on living room light"}},
		{"range of minutes", []string{"22-26 13 * * * turn on living room light"}, "04 Jan 18 13:25 UTC", []string{"turn on living room light"}},
		{"range of hours", []string{"25 13-14 * * * turn on living room light"}, "04 Jan 18 13:25 UTC", []string{"turn on living room light"}},
		{"another range of hours", []string{"25 12-14 * * * turn on living room light"}, "04 Jan 18 13:25 UTC", []string{"turn on living room light"}},
		{"tabs and space", []string{"25	13     *	*    *    turn on living room light"}, "04 Jan 18 13:25 UTC", []string{"turn on living room light"}},

		//misses
		{"nothing for minute", []string{"25 * * * * turn on living room light"}, "01 Jan 18 13:24 UTC", []string{}},
		{"nothing for minute and hour", []string{"25 13 * * * turn on living room light"}, "01 Jan 18 12:25 UTC", []string{}},
		{"nothing for minute, hour, and day", []string{"25 13 1 * * turn on living room light"}, "02 Jan 18 13:25 UTC", []string{}},
		{"nothing for minute, hour, day, month", []string{"25 13 1 1 * turn on living room light"}, "02 Jan 18 13:25 UTC", []string{}},
		{"nothing for weekday", []string{"25 13 * * 5 turn on living room light"}, "06 Jan 18 13:25 UTC", []string{}},
		{"nothing for multiple weekdays", []string{"25 13 * * 1,3,6 turn on living room light"}, "02 Jan 18 13:25 UTC", []string{}},
		{"nothing for range of weekdays", []string{"25 13 * * 4-6 turn on living room light"}, "01 Jan 18 13:25 UTC", []string{}},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%s in %v", tc.doc, tc.cmds), func(t *testing.T) {
			tm, err := time.Parse(time.RFC822, tc.time)
			g.Expect(err).To(BeNil())

			c, err := cron.New(tc.cmds)
			g.Expect(err).To(BeNil())
			out := c.Check(tm)
			g.Expect(out).To(ConsistOf(tc.want))
		})
	}
}
