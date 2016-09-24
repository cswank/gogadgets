package gogadgets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cswank/rex"
)

type Client struct {
	addr       string
	gadgetAddr string
	msg        chan Message
}

func NewClient(addr, gadgetAddr string) *Client {
	return &Client{
		addr:       addr,
		gadgetAddr: gadgetAddr,
		msg:        make(chan Message),
	}
}

func (c *Client) Connect() chan Message {
	c.register()
	r := rex.New("main")
	r.Post("/messages", http.HandlerFunc(c.update))
	go func() {
		if err := http.ListenAndServe(c.addr, r); err != nil {
			log.Fatal(err)
		}
	}()
	return c.msg
}

func (c *Client) update(w http.ResponseWriter, r *http.Request) {

	var msg Message
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&msg); err != nil {
		log.Println(err)
		return
	}
	if msg.UUID == "" {
		msg.UUID = GetUUID()
	}
	c.msg <- msg
}

func (c *Client) register() {
	var tries int
	addr := fmt.Sprintf("%s/clients", c.gadgetAddr)
	a := map[string]string{"address": fmt.Sprintf("%s/messages", c.addr), "token": "n/a"}
	for {
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.Encode(&a)
		r, err := http.Post(addr, "application/json", buf)
		if err == nil && r.StatusCode == http.StatusOK {
			return
		}
		log.Printf("unable to register, trying again: %v\n", err)
		tries = increment(tries)
		time.Sleep(time.Duration(tries) * 100 * time.Millisecond)
	}
}
