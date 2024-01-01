package gogadgets

import (
	"log"
)

type GadgetConfig struct {
	Type         string                 `json:"type,omitempty"`
	Location     string                 `json:"location,omitempty"`
	Name         string                 `json:"name,omitempty"`
	OnCommands   []string               `json:"on_commands,omitempty"`
	OffCommands  []string               `json:"off_commands,omitempty"`
	OnValue      string                 `json:"on_value,omitempty"`
	OffValue     string                 `json:"off_value,omitempty"`
	InitialValue string                 `json:"initial_value,omitempty"`
	Pin          Pin                    `json:"pin,omitempty"`
	Args         map[string]interface{} `json:"args,omitempty"`
}

type Config struct {
	Master    string         `json:"master,omitempty"`
	Host      string         `json:"host,omitempty"`
	Port      int            `json:"port,omitempty"`
	Gadgets   []GadgetConfig `json:"gadgets,omitempty"`
	Endpoints []HTTPHandler
}

func (c Config) CreateGadgets(gadgets ...Gadgeter) []Gadgeter {
	var out []Gadgeter
	out = append(out, gadgets...)
	for _, config := range c.Gadgets {
		gadget, err := NewGadget(&config)
		if err != nil {
			log.Fatal(err)
		}
		out = append(out, gadget)
	}
	out = append(out, NewServer(c.Host, c.Master, c.Port, c.Endpoints...), &MethodRunner{})
	return out
}
