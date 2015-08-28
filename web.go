package gogadgets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	host        string
	port        int
	prefix      string
	master      bool
	lg          Logger
	external    chan Message
	internal    chan Message
	id          string
	updates     map[string]Message
	seen        map[string]time.Time
	statusLock  sync.Mutex
	seenLock    sync.Mutex
	clientsLock sync.Mutex
	clients     map[string]bool
}

func NewServer(host string, port int, master bool, lg Logger) *Server {
	return &Server{
		host:     host,
		port:     port,
		master:   master,
		lg:       lg,
		updates:  map[string]Message{},
		id:       "server",
		external: make(chan Message),
		seen:     map[string]time.Time{},
		clients:  map[string]bool{},
	}
}

func (s *Server) Start(i <-chan Message, o chan<- Message) {

	go s.startServer()
	go s.cleanup()

	for {
		select {
		case msg := <-i:
			if msg.Type == UPDATE {
				s.statusLock.Lock()
				s.updates[msg.Sender] = msg
				s.statusLock.Unlock()
			}
			s.send(msg)
		case msg := <-s.external:
			s.setSeen(msg)
			o <- msg
		}
	}
}

func (s *Server) send(msg Message) {
	if s.isSeen(msg) {
		return
	}
	s.clientsLock.Lock()
	for host := range s.clients {
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		enc.Encode(msg)
		addr := fmt.Sprintf("%s:6111/gadgets", host)
		r, err := http.Post(addr, "application/json", &buf)
		if err != nil || r.StatusCode != http.StatusOK {
			delete(s.clients, host)
		}
	}
	s.clientsLock.Unlock()
}

func (s *Server) setSeen(msg Message) {
	s.seenLock.Lock()
	s.seen[msg.UUID] = time.Now()
	s.seenLock.Unlock()
}

func (s *Server) isSeen(msg Message) bool {
	s.seenLock.Lock()
	_, ok := s.seen[msg.UUID]
	if ok {
		delete(s.seen, msg.UUID)
	}
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

func (s *Server) startServer() {
	r := mux.NewRouter()
	r.HandleFunc("/gadgets", s.status).Methods("GET")
	r.HandleFunc("/gadgets", s.update).Methods("PUT", "POST")
	r.HandleFunc("/clients", s.setClient).Methods("POST")
	r.HandleFunc("/clients", s.getClients).Methods("GET")
	r.HandleFunc("/clients", s.removeClient).Methods("DELETE")

	s.lg.Printf("listening on 0.0.0.0:%d\n", s.port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), r)
	if err != nil {
		s.lg.Fatal(err)
	}
}

func (s *Server) getClients(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	s.clientsLock.Lock()
	enc.Encode(s.clients)
	s.clientsLock.Unlock()
}

func (s *Server) setClient(w http.ResponseWriter, r *http.Request) {
	var a map[string]string
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&a)
	if err != nil || len(a["address"]) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	s.clientsLock.Lock()
	s.clients[a["address"]] = true
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

func (s *Server) update(w http.ResponseWriter, r *http.Request) {
	var msg Message
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&msg); err != nil {
		lg.Println(err)
		return
	}
	if msg.UUID == "" {
		msg.UUID = getUUID()
	}
	s.external <- msg
}
