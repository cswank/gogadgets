package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/cswank/gogadgets"
)

const (
	defaultDir    = "~/.gadgets"
	defaultConfig = "/Users/Cswank/.gadgets/config.json"
)

var (
	run    = kingpin.Command("run", "run a gogadgets instance")
	config = run.Flag("config", "Path to a Gadgets config file").Short('c').Default("/etc/gogadgets.json").String()

	cmd  = kingpin.Command("cmd", "send a command to a gogadgets instance")
	host = cmd.Flag("host", "Name of gogadgets host").Short('h').Default("localhost").String()

	//rcl: robobt command language
	rcl = cmd.Arg("cmd", "a Robot Command Language string").String()

	status     = kingpin.Command("status", "get the status of a gadgets system")
	statusHost = status.Flag("host", "Name of gogadgets host").Short('h').Default("localhost").String()
	verbose    = status.Flag("verbose", "get the verbose status of a gadgets system").Short('v').Bool()
)

func main() {
	kingpin.Version(gogadgets.Version)
	switch kingpin.Parse() {
	case "run":
		runGadgets()
	case "cmd":
		sendCommand()
	case "status":
		getStatus()
	default:
		log.Fatal("unknown command")
	}
}

func runGadgets() {
	gogadgets.New(getConfig()).Start()
}

func getConfig() string {
	cfg := *config
	if cfg != "" {
		return cfg
	}
	if gogadgets.FileExists(defaultConfig) {
		return defaultConfig
	}
	return ""
}

func getStatus() {
	addr := fmt.Sprintf("http://%s:%d/gadgets", *statusHost, 6111)
	if *verbose {
		getVerbose(addr)
		return
	}

	r, err := http.Get(fmt.Sprintf("%s/values", addr))
	if err != nil {
		log.Fatal("err", err)
	}

	var v map[string]map[string]gogadgets.Value
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&v); err != nil {
		log.Fatal("err", err)
	}

	d, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(d))
}

func getVerbose(addr string) {
	r, err := http.Get(addr)
	if err != nil {
		log.Fatal("err", err)
	}

	var s map[string]gogadgets.Message
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&s); err != nil {
		log.Fatal("err", err)
	}

	d, _ := json.MarshalIndent(s, "", "  ")
	fmt.Println(string(d))
}

func sendCommand() {
	msg := gogadgets.Message{
		UUID:   gogadgets.GetUUID(),
		Type:   gogadgets.COMMAND,
		Sender: "client",
		Body:   *rcl,
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.Encode(&msg)
	addr := fmt.Sprintf("http://%s:%d/gadgets", *host, 6111)
	r, err := http.Post(addr, "application/json", &buf)
	if err != nil {
		log.Fatal("err", err)
	}
	fmt.Println(r.Status)
}
