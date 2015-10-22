package gogadgets

import (
	"fmt"
	"log"

	"github.com/brutella/hc/hap"
	"github.com/brutella/hc/model"
	"github.com/brutella/hc/model/accessory"
)

type HomeKit struct {
	id          string
	switches    map[string]model.Switch
	accessories []*accessory.Accessory
	key         string
}

func NewHomeKit(key string, gadgets []Gadgeter) *HomeKit {
	s, a := getSwitches(gadgets)
	return &HomeKit{
		id:          "homekit",
		switches:    s,
		accessories: a,
		key:         key,
	}
}

func getSwitches(gadgets []Gadgeter) (map[string]model.Switch, []*accessory.Accessory) {
	switches := map[string]model.Switch{}
	var accessories []*accessory.Accessory
	for _, g := range gadgets {
		if g.GetDirection() != "output" {
			info := model.Info{
				Name:         g.GetUID(),
				Manufacturer: "gogadgets",
			}
			s := accessory.NewSwitch(info)
			switches[g.GetUID()] = s
			accessories = append(accessories, s.Accessory)
		}
	}
	return switches, accessories
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
	var t hap.Transport
	var err error
	if len(h.accessories) > 0 {
		t, err = hap.NewIPTransport(h.key, h.accessories[0], h.accessories[1:]...)
	} else {
		t, err = hap.NewIPTransport(h.key, h.accessories[0])
	}
	if err != nil {
		log.Fatal(err)
	}
	t.Start()
}

func (h *HomeKit) GetUID() string {
	return h.id
}

func (h *HomeKit) GetDirection() string {
	return "na"
}
