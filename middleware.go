package main

import (
	"log"
	"net/http"
	"time"
)

type Node struct {
	Path       string
	Dispatcher http.HandlerFunc
}

type Middleware struct {
	Nodes    map[string][]*Node
	NotFound http.HandlerFunc
}

func NewMiddleware() *Middleware {
	return &Middleware{Nodes: make(map[string][]*Node)}
}

func (m *Middleware) Dispatcher(w http.ResponseWriter, r *http.Request) {
	children, ok := m.Nodes[r.Method]

	if !ok {
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	for _, child := range children {
		if child.Path == r.URL.Path {
			child.Dispatcher(w, r)
			return
		}
	}

	http.NotFound(w, r)
}

func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var query string
	var start = time.Now()
	var writer = StatusWriter{w, 0, 0}

	if r.URL.RawQuery != "" {
		query = "?" + r.URL.RawQuery
	}

	m.Dispatcher(&writer, r)

	log.Printf("%s %s \"%s %s %s\" %d %d \"%s\" %v",
		r.Host,
		r.RemoteAddr,
		r.Method,
		r.URL.Path+query,
		r.Proto,
		writer.Status,
		writer.Length,
		r.Header.Get("User-Agent"),
		time.Now().Sub(start))
}

func (m *Middleware) Handle(method, path string, handle http.HandlerFunc) {
	m.Nodes[method] = append(m.Nodes[method], &Node{path, handle})
}

func (m *Middleware) GET(path string, handle http.HandlerFunc) {
	m.Handle("GET", path, handle)
}

func (m *Middleware) POST(path string, handle http.HandlerFunc) {
	m.Handle("POST", path, handle)
}
