package rex

import (
	"net/http"
	"strings"
)

var (
	routers map[string]*Router
)

func init() {
	routers = map[string]*Router{}
}

func New(name string) *Router {
	r := &Router{
		methods: map[string]*node{
			"GET":    newNode(),
			"POST":   newNode(),
			"PUT":    newNode(),
			"DELETE": newNode(),
			"PATCH":  newNode(),
			"HEAD":   newNode(),
		},
	}
	routers[name] = r
	return r
}

type Router struct {
	fileServer http.Handler
	methods    map[string]*node
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	n, ok := r.methods[req.Method]
	if ok {
		found := n.handle(strings.Split(strings.Trim(req.URL.Path, "/"), "/"), w, req)
		if !found && r.fileServer != nil {
			r.fileServer.ServeHTTP(w, req)
		} else if !found {
			w.WriteHeader(http.StatusNotFound)
		}
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func Vars(r *http.Request, name string) map[string]string {
	pth := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	m := map[string]string{}
	t, ok := routers[name]
	if !ok {
		return m
	}
	t.methods[r.Method].vars(pth, m)
	return m
}

func (r *Router) ServeFiles(h http.Handler) {
	r.fileServer = h
}

func (r *Router) Get(pth string, h http.Handler) *node {
	return r.method("GET", pth, h)
}

func (r *Router) Post(pth string, h http.Handler) *node {
	return r.method("POST", pth, h)
}

func (r *Router) Put(pth string, h http.Handler) *node {
	return r.method("PUT", pth, h)
}

func (r *Router) Delete(pth string, h http.Handler) *node {
	return r.method("DELETE", pth, h)
}

func (r *Router) Path(pth string, h http.Handler) *node {
	return r.method("PATCH", pth, h)
}

func (r *Router) method(name, pth string, h http.Handler) *node {
	return r.methods[name].add(strings.Split(strings.Trim(pth, "/"), "/"), h)
}

type node struct {
	id       string
	key      string
	handler  http.Handler
	resource bool
	children map[string]*node
}

func newNode() *node {
	return &node{
		children: map[string]*node{},
	}
}

func (n *node) vars(pth []string, m map[string]string) {
	if len(pth) == 0 && !n.resource {
		return
	}

	var c *node
	var ok bool
	l := len(pth)
	if n.resource && l != 0 {
		m[n.key] = pth[0]
		c, ok = n.children[n.id]
	} else if l != 0 {
		c, ok = n.children[pth[0]]
	}

	if ok && l != 0 {
		c.vars(pth[1:], m)
	}
}

func (n *node) add(pth []string, h http.Handler) *node {
	if len(pth) == 0 {
		n.handler = h
		return n
	}
	n.resource, n.id, n.key = n.isResource(pth[0])
	var x string
	if n.resource {
		x = n.id
	} else {
		x = pth[0]
	}
	c, ok := n.children[x]
	if !ok {
		c = newNode()
		n.children[pth[0]] = c
	}
	return c.add(pth[1:], h)
}

func (n *node) isResource(s string) (bool, string, string) {
	var id string
	r := strings.Index(s, "{") == 0
	if r {
		id = s
	}
	return r, id, strings.Trim(id, "{}")
}

func (n *node) handle(pth []string, w http.ResponseWriter, r *http.Request) bool {
	if len(pth) == 0 {
		if n == nil || n.handler == nil {
			return false
		}
		n.handler.ServeHTTP(w, r)
		return true
	}
	var x string
	if n.resource {
		x = n.id
	} else {
		x = pth[0]
	}
	if c, ok := n.children[x]; ok {
		return c.handle(pth[1:], w, r)
	}
	return false
}
