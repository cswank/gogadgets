package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/cswank/gogadgets"
	"github.com/cswank/gogadgets/utils"
)

const (
	defaultDir    = "~/.gadgets"
	defaultConfig = "/Users/Cswank/.gadgets/config.json"
)

var (
	host    = flag.String("h", "localhost", "Name of Host")
	config  = flag.String("g", "", "Path to a Gadgets config file")
	cmd     = flag.String("c", "", "a Robot Command Language string")
	status  = flag.Bool("s", false, "get the status of a gadgets system")
	verbose = flag.Bool("v", false, "get the status of a gadgets system")
	addr    string
)

func main() {
	flag.Parse()
	addr = fmt.Sprintf("http://%s:%d/gadgets", *host, 6111)
	if len(*cmd) > 0 {
		sendCommand()
	} else if *status {
		getStatus()
	} else if *verbose {
		getVerbose()
	} else {
		runGadgets()
	}
}

func runGadgets() {
	cfg := getConfig()
	if cfg == "" {
		listen()
	} else {
		a := gogadgets.NewApp(cfg)
		a.Start()
	}
}

func getConfig() string {
	cfg := *config
	if cfg != "" {
		return cfg
	}
	if utils.FileExists(defaultConfig) {
		return defaultConfig
	}
	return ""
}

func getStatus() {

	r, err := http.Get(addr)
	if err != nil {
		log.Fatal("err", err)
	}
	var s map[string]gogadgets.Message
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&s); err != nil {
		log.Fatal("err", err)
	}

	v := map[string]map[string]gogadgets.Value{}

	for _, msg := range s {
		l, ok := v[msg.Location]
		if !ok {
			l = map[string]gogadgets.Value{}
		}
		l[msg.Name] = msg.Value
		v[msg.Location] = l
	}

	d, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(d))
}

func getVerbose() {

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
		Body:   *cmd,
	}
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.Encode(&msg)
	r, err := http.Post(addr, "application/json", &buf)
	if err != nil {
		log.Fatal("err", err)
	}
	fmt.Println(r.Status)
}

//Waits for a zmq message that contains a gogadgets
//config.  When one is recieved it is written to the
//default config path and a a gogadgts system is started.
func listen() {
	// cfg := gogadgets.SocketsConfig{
	// 	Host:    *host,
	// 	SubPort: 6111,
	// 	PubPort: 6112,
	// 	Master:  false,
	// }
	// s := gogadgets.NewSockets(cfg)
	// err := s.Connect()
	// if err != nil {
	// 	panic(err)
	// }
	// defer s.Close()
	// time.Sleep(100 * time.Millisecond)
	// log.Println("listening for new gadgets")
	// msg := s.Recv()
	// d, err := json.Marshal(&msg.Config)
	// if err != nil {
	// 	panic(err)
	// }
	// os.Mkdir(defaultDir, 0644)
	// err = ioutil.WriteFile(defaultConfig, d, 0644)
	// if err != nil {
	// 	panic(err)
	// }
	// time.Sleep(100 * time.Millisecond)
	// runGadgets()
}
