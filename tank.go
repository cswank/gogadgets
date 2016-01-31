package gogadgets

import "fmt"

//Watcher checks the incoming message and then returns a change
//in volume.
type watcher func(float64, Message) (float64, bool)

type source struct {
	Type string `json:"type"`
	Name string `json:"name"`
}

//Tank represents a tank that holds liquid.  It keeps track
//of how much is inside.  A tank is responsible for keeping
//track of how much liquid it has when it drains.  The downstream
//tank that this tank drains into then knows how much liquid it
//has by keeping an eye on its source (tank or valve).
type Tank struct {
	volume float64
	out    chan<- Value
	//When the sink of this tank is commanded to be filled
	//a goroutine is started that sends updates via the change
	//chan.
	change chan float64
	units  string
	//Source is the tank (or valve) that supplies
	//liquid to the tank.  The tank should update its volume
	//based on changes in the source.
	source watcher
}

func NewTank(pin *Pin) (InputDevice, error) {
	w, err := newWatcher(pin.Args)
	if err != nil {
		return nil, err
	}
	return &Tank{
		source: w,
		change: make(chan float64),
		units:  pin.Units,
	}, nil
}

func (t *Tank) SendValue() {
	t.out <- Value{
		Value: t.volume,
		Units: t.units,
	}
}

func (t *Tank) GetValue() *Value {
	return &Value{
		Value: t.volume,
		Units: t.units,
	}
}

func (t *Tank) Start(in <-chan Message, out chan<- Value) {
	t.out = out
	keepGoing := true
	for keepGoing {
		select {
		case msg := <-in:
			t.update(msg)
		case v := <-t.change:
			t.volume = v
			t.SendValue()
		}
	}
}

func (t *Tank) update(msg Message) {
	v, changed := t.source(t.volume, msg)
	if changed {
		t.volume = v
		t.SendValue()
	}
}

func (t *Tank) Config() ConfigHelper {
	return ConfigHelper{}
}

func newWatcher(args map[string]interface{}) (watcher, error) {
	switch args["type"].(string) {
	case "valve":
		return newValveWatcher(args["float_switch"].(string), args["full"].(float64))
	case "tank":
		return newTankWatcher(args["source"].(string))
	}
	return nil, fmt.Errorf("source.type must be either 'valve' or 'tank'")
}

func newValveWatcher(name string, full float64) (watcher, error) {
	return func(vol float64, msg Message) (float64, bool) {
		if msg.Sender == name && msg.Value.Value == true {
			return full, true
		}
		return 0.0, false
	}, nil
}

func newTankWatcher(name string) (watcher, error) {
	return func(vol float64, msg Message) (float64, bool) {
		d := msg.Value.Diff
		if msg.Sender == name && d != 0.0 {
			return vol + d, true
		}
		return 0.0, false
	}, nil
}
