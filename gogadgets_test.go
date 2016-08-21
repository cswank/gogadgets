package gogadgets_test

import (
	"time"

	"github.com/cswank/gogadgets"
)

type FakeOutput struct {
	on bool
}

func (f *FakeOutput) Commands(l, n string) *gogadgets.Commands {
	return nil
}

func (f *FakeOutput) Config() gogadgets.ConfigHelper {
	return gogadgets.ConfigHelper{}
}

func (f *FakeOutput) Update(msg *gogadgets.Message) bool {
	return false
}

func (f *FakeOutput) On(val *gogadgets.Value) error {
	f.on = true
	return nil
}

func (f *FakeOutput) Off() error {
	f.on = false
	return nil
}

func (f *FakeOutput) Status() map[string]bool {
	return map[string]bool{"gpio": f.on}
}

type FakePoller struct {
	trigger chan bool
	val     bool
}

func (f *FakePoller) Status() map[string]bool {
	return map[string]bool{"poller": f.val}
}

func (f *FakePoller) Wait() (bool, error) {
	if f.trigger == nil {
		time.Sleep(100 * time.Millisecond)
		f.val = !f.val
	} else {
		f.val = <-f.trigger
	}
	return f.val, nil
}
