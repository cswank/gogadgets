package gogadgets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/cswank/rex"
)

type Server struct {
	host        string
	master      string
	isMaster    bool
	port        int
	prefix      string
	lg          Logger
	external    chan Message
	internal    chan Message
	id          string
	updates     map[string]Message
	seen        map[string]time.Time
	statusLock  sync.Mutex
	seenLock    sync.Mutex
	clientsLock sync.Mutex
	clients     map[string]string
}

func NewServer(host, master string, port int, lg Logger) *Server {
	var isMaster bool
	clients := map[string]string{}
	if master == "" {
		isMaster = true
	} else {
		clients = map[string]string{
			fmt.Sprintf("%s/gadgets", master): "master",
		}
	}
	return &Server{
		master:   master,
		host:     host,
		isMaster: isMaster,
		port:     port,
		lg:       lg,
		updates:  map[string]Message{},
		id:       "server",
		external: make(chan Message),
		seen:     map[string]time.Time{},
		clients:  clients,
	}
}

func (s *Server) Start(i <-chan Message, o chan<- Message) {
	if !s.isMaster {
		s.register()
	}
	go s.startServer()
	go s.cleanup()

	for {
		select {
		case msg := <-i:
			if (msg.Type == UPDATE || msg.Type == METHODUPDATE) && s.isMaster {
				s.statusLock.Lock()
				s.updates[msg.Sender] = msg
				s.statusLock.Unlock()
			}
			if !s.isSeen(msg) {
				s.send(msg)
			}
		case msg := <-s.external:
			s.setSeen(msg)
			o <- msg
			if s.isMaster {
				s.send(msg)
			}
		}
	}
}

func (s *Server) send(msg Message) {
	s.clientsLock.Lock()
	for host, token := range s.clients {
		go s.doSend(host, msg, token)
	}
	s.clientsLock.Unlock()
}

func (s *Server) doSend(host string, msg Message, token string) {
	if len(msg.Host) > 0 && strings.Index(host, msg.Host) == 0 {
		return
	}
	msg.Host = s.host
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.Encode(msg)
	req, err := http.NewRequest("POST", host, &buf)

	s.clientsLock.Lock()
	defer s.clientsLock.Unlock()
	if err != nil {
		delete(s.clients, host)
		log.Printf("error posting to quimby host: %s, err: %v\n", host, err)
		return
	}
	if len(token) > 0 {
		req.Header.Add("Authorization", token)
	}
	r, err := http.DefaultClient.Do(req)
	if r != nil {
		r.Body.Close()
	}
	if err != nil || r.StatusCode != http.StatusOK && token != "master" {
		delete(s.clients, host)
	}
}

func (s *Server) setSeen(msg Message) {
	s.seenLock.Lock()
	s.seen[msg.UUID] = time.Now()
	s.seenLock.Unlock()
}

func (s *Server) isSeen(msg Message) bool {
	s.seenLock.Lock()
	_, ok := s.seen[msg.UUID]
	s.seenLock.Unlock()
	return ok
}

func (s *Server) cleanup() {
	for {
		time.Sleep(60 * time.Second)
		now := time.Now()
		s.seenLock.Lock()
		for k, v := range s.seen {
			if v.Sub(now) > 10*time.Second {
				delete(s.seen, k)
			}
		}
		s.seenLock.Unlock()
	}
}

func (s *Server) GetUID() string {
	return s.id
}

func (s *Server) GetDirection() string {
	return "na"
}

func (s *Server) startServer() {
	r := rex.New("main")
	r.Get("/gadgets", http.HandlerFunc(s.status))
	r.Put("/gadgets", http.HandlerFunc(s.update))
	r.Post("/gadgets", http.HandlerFunc(s.update))
	r.Get("/gadgets/values", http.HandlerFunc(s.values))
	r.Get("/gadgets/locations/{location}/devices/{device}/status", http.HandlerFunc(s.deviceValue))
	if s.isMaster {
		r.Post("/clients", http.HandlerFunc(s.setClient))
		r.Get("/clients", http.HandlerFunc(s.getClients))
		r.Delete("/clients", http.HandlerFunc(s.removeClient))
	}

	s.lg.Printf("listening on 0.0.0.0:%d\n", s.port)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {
		s.lg.Fatal(err)
	}
}

func (s *Server) getClients(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	s.clientsLock.Lock()
	if err := enc.Encode(s.clients); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	s.clientsLock.Unlock()
}

func (s *Server) setClient(w http.ResponseWriter, r *http.Request) {
	var a map[string]string
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&a)
	if err != nil || len(a["address"]) == 0 || len(a["token"]) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	s.clientsLock.Lock()
	s.clients[a["address"]] = a["token"]
	s.clientsLock.Unlock()
}

func (s *Server) removeClient(w http.ResponseWriter, r *http.Request) {
	s.clientsLock.Lock()
	delete(s.clients, r.URL.Host)
	s.clientsLock.Unlock()
}

func (s *Server) status(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	s.statusLock.Lock()
	if err := enc.Encode(s.updates); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		lg.Println(err)
	}
	s.statusLock.Unlock()
}

func (s *Server) deviceValue(w http.ResponseWriter, r *http.Request) {
	vars := rex.Vars(r, "main")
	k := fmt.Sprintf("%s %s", vars["location"], vars["device"])
	s.statusLock.Lock()
	m, ok := s.updates[k]
	s.statusLock.Unlock()
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	enc := json.NewEncoder(w)
	if err := enc.Encode(m.Value); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *Server) values(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	s.statusLock.Lock()
	v := map[string]map[string]Value{}

	for _, msg := range s.updates {
		l, ok := v[msg.Location]
		if !ok {
			l = map[string]Value{}
		}
		l[msg.Name] = msg.Value
		v[msg.Location] = l
	}
	s.statusLock.Unlock()
	if err := enc.Encode(v); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		lg.Println(err)
	}
}

func (s *Server) update(w http.ResponseWriter, r *http.Request) {
	var msg Message
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&msg); err != nil {
		lg.Println(err)
		return
	}
	if msg.UUID == "" {
		msg.UUID = GetUUID()
	}
	s.external <- msg
}

//
func (s *Server) register() {
	var tries int
	addr := fmt.Sprintf("%s/clients", s.master)
	a := map[string]string{"address": fmt.Sprintf("%s/gadgets", s.host), "token": "n/a"}
	for {
		buf := &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		enc.Encode(&a)
		r, err := http.Post(addr, "application/json", buf)
		if err == nil && r.StatusCode == http.StatusOK {
			return
		}
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
