package gogadgets

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"
)

type FakeOutput struct {
	OutputDevice
	on bool
}

func (f *FakeOutput) Update(msg *Message) {

}

func (f *FakeOutput) On(val *Value) error {
	f.on = true
	return nil
}

func (f *FakeOutput) Off() error {
	f.on = false
	return nil
}

func (f *FakeOutput) Status() interface{} {
	return f.on
}

type FakePoller struct {
	Poller
	val bool
}

func (f *FakePoller) Wait() (bool, error) {
	time.Sleep(100 * time.Millisecond)
	f.val = !f.val
	return f.val, nil
}

func TestGadgets(t *testing.T) {
	port := 1024 + rand.Intn(65535-1024)
	p := &Gadget{
		Location:   "tank",
		Name:       "pump",
		OnCommand:  fmt.Sprintf("turn on %s %s", "tank", "pump"),
		OffCommand: fmt.Sprintf("turn off %s %s", "tank", "pump"),
		Output:     &FakeOutput{},
		UID:        fmt.Sprintf("%s %s", "tank", "pump"),
	}
	location := "tank"
	name := "switch"
	poller := &FakePoller{}
	s := &Gadget{
		Location: location,
		Name:     name,
		Input: &Switch{
			GPIO:  poller,
			Value: 5.0,
			TrueValue: 5.0,
			Units: "liters",
		},
		UID: fmt.Sprintf("%s %s", location, name),
	}
	a := App{
		Gadgets: []GoGadget{p, s},
		Host:    "localhost",
		SubPort: port,
		PubPort: port + 1,
	}
	go a.Start()
	time.Sleep(1 * time.Second)
}

func TestGetConfigFromFile(t *testing.T) {
	cfg := `{
    "gadgets": [
        {
            "location": "front yard",
            "name": "sprinklers",
            "pin": {
                "type": "gpio",
                "port": "8",
                "pin": "10",
                "direction": "out"
            }
        },
        {
            "location": "front garden",
            "name": "sprinklers",
            "pin": {
                "type": "gpio",
                "port": "8",
                "pin": "11",
                "direction": "out"
            }
        },
        {
            "location": "sidewalk",
            "name": "sprinklers",
            "pin": {
                "type": "gpio",
                "port": "8",
                "pin": "12",
                "direction": "out"
            }
        },
        {
            "location": "back yard",
            "name": "sprinklers",
            "pin": {
                "type": "gpio",
                "port": "8",
                "pin": "14",
                "direction": "out"
            }
        },
        {
            "location": "back garden",
            "name": "sprinklers",
            "pin": {
                "type": "gpio",
                "port": "8",
                "pin": "15",
                "direction": "out"
            }
        }
    ]
}
`
	f, _ := ioutil.TempFile("", "")
	f.Write([]byte(cfg))
	f.Close()
	config := getConfig(f.Name())
	os.Remove(f.Name())
	if len(config.Gadgets) != 5 {
		t.Error(config)
	}
}
