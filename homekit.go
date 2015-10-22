package gogadgets

import (
	"fmt"

	"github.com/brutella/hc/model"
	"github.com/brutella/hc/model/accessory"
)

type HomeKit struct {
	id       string
	switches map[string]model.Switch
	key      string
}

func NewHomeKit(gadgets []Gadgeter) *HomeKit {
	return &HomeKit{
		id:       "homekit",
		switches: getSwitches(gadgets),
		key:      "44444444",
	}
}

func getSwitches(gadgets []Gadgeter) map[string]model.Switch {
	switches := map[string]model.Switch{}
	for _, g := range gadgets {
		if g.GetDirection() != "output" {
			info := model.Info{
				Name:         g.GetUID(),
				Manufacturer: "gogadgets",
			}
			s := accessory.NewSwitch(info)
			switches[g.GetUID()] = s
		}
	}
	return switches
}

func (h *HomeKit) Start(i <-chan Message, o chan<- Message) {
	for k, s := range h.switches {
		s.OnStateChanged(func(on bool) {
			if on == true {
				o <- Message{
					Type: COMMAND,
					Body: fmt.Sprintf("turn on %s", k),
				}
			} else {
				o <- Message{
					Type: COMMAND,
					Body: fmt.Sprintf("turn off %s", k),
				}
			}
		})
	}
	for {
		select {
		case <-i:
		}
	}
}

func (h *HomeKit) GetUID() string {
	return h.id
}

func (h *HomeKit) GetDirection() string {
	return "na"
}
