Rex
===

The RESTful mux package.  If you have a RESTful
set or routes like:

    /
    /friends
    /friends/{friend}
    /friends/{friend}/hobbies/{hobby}
    /friends/{friend}/songs/{song}

Then you should be able to use Rex.  It doesn't match
routes based on regex and only provides enough functionality
to work with a RESTful api.  Rex can be used to serve the above
routes like this:

    package main

    import (
    	"log"
    	"net/http"

    	"github.com/cswank/rex"
    )

    func main() {
    	r := rex.New("example")
    	r.Get("/things", getThings)
        r.Post("/things", addThing)
    	r.Get("/things/{thing}", getThing)
        r.Delete("/things/{thing}", deleteThing)

    	http.Handle("/", r)
    	log.Println(http.ListenAndServe(":8888", r))
    }

    func getThings(w http.ResponseWriter, r *http.Request) {
    	w.Write([]byte("things"))
    }

    func addThing(w http.ResponseWriter, r *http.Request) {
        
    }

    func getThing(w http.ResponseWriter, r *http.Request) {
    	w.Write([]byte("thing"))
    }

    func deleteThing(w http.ResponseWriter, r *http.Request) {
    
    }

You can also serve files:

    package main

    import (
    	"log"
    	"net/http"

    	"github.com/cswank/rex"
    )

    func main() {
    	r := rex.New("example")
    	r.Get("/things", getThings)
        r.ServeFiles(http.FileServer(http.Dir("/www")))
    	http.Handle("/", r)
    	log.Println(http.ListenAndServe(":8888", r))
    }

    func getThings(w http.ResponseWriter, r *http.Request) {
    	w.Write([]byte("things"))
    }

The supported HTTP methods are Get, Post, Put, Delete, Patch, and Head.

When you create a new Router you must pass in a name.  When getting
the vars from a request the name of the router must be passed in.


    func getThings(w http.ResponseWriter, r *http.Request) {
        vars := rex.Vars("example")
    	w.Write([]byte("things"))
    }


Compared to gorilla/mux it seems to be a bit faster when you
run the benchmark tests.  Running the rex benchmark:

     package rex_test

     import (
     	"net/http"
     	"testing"

     	"github.com/cswank/rex"
     )

     func BenchmarkMux(b *testing.B) {
     	r := rex.New("bench")
     	handler := func(w http.ResponseWriter, r *http.Request) {}
     	r.Get("/v1/{v1}", handler)

     	request, _ := http.NewRequest("GET", "/v1/anything", nil)
     	for i := 0; i < b.N; i++ {
     		r.ServeHTTP(nil, request)
     	}
     }    

Gives

    ➜  rex git:(master) go test -bench=.
    BenchmarkRex-2	 5000000	       370 ns/op
    ok  	github.com/cswank/rex	2.259s

And running the gorilla/mux benchmark test:


    // Copyright 2012 The Gorilla Authors. All rights reserved.
    // Use of this source code is governed by a BSD-style
    // license that can be found in the LICENSE file.

    package mux

    import (
    	"net/http"
    	"testing"
    )

    func BenchmarkMux(b *testing.B) {
    	router := new(Router)
    	handler := func(w http.ResponseWriter, r *http.Request) {}
    	router.HandleFunc("/v1/{v1}", handler)

    	request, _ := http.NewRequest("GET", "/v1/anything", nil)
    	for i := 0; i < b.N; i++ {
    		router.ServeHTTP(nil, request)
    	}
    }

Gives:
    
    ➜  mux git:(f15e0c4) go test -bench=.
    PASS
    BenchmarkMux-2	  500000	      2728 ns/op
    ok  	github.com/gorilla/mux	1.418s
  
On the same machine.
