package gogadgets

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type Server struct {
	host     string
	port     int
	master   bool
	lg       Logger
	external chan Message
	internal chan Message
	id       string
	updates  map[string]Message
	lock     sync.Mutex
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
	}
}

func (s *Server) Start(i <-chan Message, o chan<- Message) {

	go s.startServer()

	for {
		select {
		case msg := <-i:
			if msg.Type == UPDATE {
				s.lock.Lock()
				s.updates[msg.Sender] = msg
				s.lock.Unlock()
			}
		case msg := <-s.external:
			o <- msg
		}
	}
}

func (s *Server) GetUID() string {
	return s.id
}

func (s *Server) startServer() {
	r := mux.NewRouter()
	r.HandleFunc("/gadgets", s.status).Methods("GET")
	r.HandleFunc("/gadgets", s.update).Methods("PUT", "POST")

	http.Handle("/", r)

	s.lg.Printf("listening on 0.0.0.0:%d/n", s.port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", s.port), r)
	if err != nil {
		s.lg.Fatal(err)
	}
}

func (s *Server) status(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	s.lock.Lock()
	if err := enc.Encode(s.updates); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		lg.Println(err)
	}
	s.lock.Unlock()
}

func (s *Server) update(w http.ResponseWriter, r *http.Request) {
	var msg Message
	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&msg); err != nil {
		lg.Println(err)
		return
	}
	s.external <- msg
}
