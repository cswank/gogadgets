package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cswank/gogadgets"
	"github.com/cswank/rex"
)

type Client struct {
	addr       string
	gadgetAddr string
	msg        chan gogadgets.Message
}

func New(addr, gadgetAddr string) *Client {
	return &Client{
		addr:       addr,
		gadgetAddr: gadgetAddr,
		msg:        make(chan gogadgets.Message),
	}
}

func (c *Client) Connect() chan gogadgets.Message {
	c.register()
	r := rex.New("main")
	r.Post("/messages", http.HandlerFunc(c.update))
	r.Get("/ping", http.HandlerFunc(c.ping))
	go func() {
		if err := http.ListenAndServe(c.addr, r); err != nil {
			log.Fatal(err)
		}
	}()

	for {
		resp, err := http.Get(fmt.Sprintf("http://%s/ping", c.addr))
		if err != nil {
			time.Sleep(10 * time.Millisecond)
		} else {
			resp.Body.Close()
			break
		}
	}
	return c.msg
}

func (c *Client) ping(w http.ResponseWriter, r *http.Request) {
	r.Body.Close()
}

func (c *Client) update(w http.ResponseWriter, r *http.Request) {

	var msg gogadgets.Message
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&msg); err != nil {
		log.Println(err)
		return
	}
	r.Body.Close()
	if msg.UUID == "" {
		msg.UUID = gogadgets.GetUUID()
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

func increment(i int) int {
	if i == 100 {
		return i
	}
	return i + 1
}
