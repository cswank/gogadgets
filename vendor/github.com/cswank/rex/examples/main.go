package main

import (
	"log"
	"net/http"

	"github.com/cswank/rex"
)

func main() {
	r := rex.New("example")
	r.Get("/things", getThings)
	r.Get("/things/{thing}", getThing)
	r.ServeFiles(http.FileServer(http.Dir("./www")))

	http.Handle("/", r)

	log.Println(http.ListenAndServe(":8888", r))
}

func getThings(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("things"))
}

func getThing(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("thing"))
}
