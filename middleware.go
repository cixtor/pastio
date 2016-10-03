package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"
)

type Node struct {
	Path            string
	Params          []string
	NumParams       int
	NumSections     int
	Dispatcher      http.HandlerFunc
	MatchEverything bool
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
		/* Internal server error if HTTP method is not allowed */
		http.Error(w, http.StatusText(405), http.StatusMethodNotAllowed)
		return
	}

	if r.URL.Path == "" || r.URL.Path[0] != '/' {
		/* Bad request error if URL does not starts with slash */
		http.Error(w, http.StatusText(400), http.StatusBadRequest)
		return
	}

	var ctx = r.Context()
	var parameters []string
	var sections []string

	var lendef int   // Length defined URL
	var lenreq int   // Length requested URL
	var extra string // Dynamic URL parameters

	for _, child := range children {
		/* If URL matches and there are no dynamic parameters */
		if child.Path == r.URL.Path && child.Params == nil {
			child.Dispatcher(w, r)
			return
		}

		/* Continue only if the defined URL contains dynamic params. */
		if child.Params == nil {
			continue
		}

		/**
		 * If the defined URL contains dynamic parameters we need to check if
		 * the requested URL is longer than the defined URL without the dynamic
		 * sections, this means that the requested URL must be longer than the
		 * clean defined URL.
		 *
		 * Defined (Raw):   /lorem/ipsum/dolor/:unique
		 * Defined (Clean): /lorem/ipsum/dolor
		 * Req. URL (Bad):  /lorem/ipsum/dolor
		 * Req. URL (Semi): /lorem/ipsum/dolor/
		 * Req. URL (Good): /lorem/ipsum/dolor/something
		 *
		 * Notice how the good requested URL has more characters than the clean
		 * defined URL, the extra characters will be extracted and converted
		 * into variables to be passed to the handler. The bad requested URL
		 * matches the exact same clean defined URL but has no extra characters,
		 * so variable "unique" will be empty which is non-processable. The semi
		 * good requested URL contains one character more than the clean defined
		 * URL, the extra character is simply a forward slash, which means the
		 * dynamic variable "unique" will be empty but at least it was on purpose.
		 *
		 * The requested URL must contains the same characters than the clean
		 * defined URL, at least from index zero, the rest of the requested URL
		 * can be different. This is to prevent malicious requests with semi
		 * valid URLs with different roots which might translate to handlers
		 * processing unrelated requests.
		 */
		lendef = len(child.Path)
		lenreq = len(r.URL.Path)
		if lendef >= lenreq {
			continue
		}

		/* Skip if root section of requested URL does not matches */
		if child.Path != r.URL.Path[0:lendef] {
			continue
		}

		/* Handle request for static files */
		if child.MatchEverything {
			child.Dispatcher(w, r)
			return
		}

		/* Separate dynamic characters from URL */
		parameters = []string{}
		extra = r.URL.Path[lendef:lenreq]
		sections = strings.Split(extra, "/")

		for _, param := range sections {
			if param != "" {
				parameters = append(parameters, param)
			}
		}

		/* Skip if number of dynamic parameters is different */
		if child.NumParams != len(parameters) {
			continue
		}

		for key, name := range child.Params {
			ctx = context.WithValue(ctx, name, parameters[key])
		}

		child.Dispatcher(w, r.WithContext(ctx))
		return
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

func (m *Middleware) ServeFiles(root string, prefix string) http.HandlerFunc {
	fs := http.FileServer(http.Dir(root))
	handler := http.StripPrefix(prefix, fs)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path[len(r.URL.Path)-1] == '/' {
			http.Error(w, http.StatusText(403), http.StatusForbidden)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func (m *Middleware) Handle(method, path string, handle http.HandlerFunc) {
	var node Node
	var parts []string
	var usable []string

	node.Path = "/"
	node.Dispatcher = handle
	parts = strings.Split(path, "/")

	// Separate dynamic parameters from the static URL.
	for _, section := range parts {
		if section == "" {
			continue
		}

		if len(section) > 1 && section[0] == ':' {
			node.Params = append(node.Params, section[1:])
			node.NumSections += 1
			node.NumParams += 1
			continue
		}

		usable = append(usable, section)
		node.NumSections += 1
	}

	node.Path += strings.Join(usable, "/")

	m.Nodes[method] = append(m.Nodes[method], &node)
}

func (m *Middleware) STATIC(root string, prefix string) {
	var node Node

	node.Path = prefix
	node.MatchEverything = true
	node.Params = []string{"filepath"}
	node.Dispatcher = m.ServeFiles(root, prefix)

	m.Nodes["GET"] = append(m.Nodes["GET"], &node)
	m.Nodes["POST"] = append(m.Nodes["POST"], &node)
}

func (m *Middleware) GET(path string, handle http.HandlerFunc) {
	m.Handle("GET", path, handle)
}

func (m *Middleware) POST(path string, handle http.HandlerFunc) {
	m.Handle("POST", path, handle)
}
