package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Pastio struct{}

func (p *Pastio) Port() string {
	var port string = os.Getenv("PORT")

	if port == "" {
		return "8080"
	}

	return port
}

func render(name string, w http.ResponseWriter) {
	tpl, err := template.ParseFiles(name)

	if err != nil {
		log.Print(err)
	} else {
		tpl.Execute(w, nil)
	}
}

func noDirectoryListing(handler http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path[len(r.URL.Path)-1] == '/' {
			http.NotFound(w, r)
			return
		}

		handler.ServeHTTP(w, r)
	})
}

func serveStaticFile(dir string) http.HandlerFunc {
	handler := http.FileServer(http.Dir(dir))
	fs := http.StripPrefix("/assets/", handler)

	return noDirectoryListing(fs)
}

func index(w http.ResponseWriter, r *http.Request) {
	render("_views/index.tmpl", w)
}

func main() {
	var pastio Pastio

	var port string = fmt.Sprintf(":%s", pastio.Port())

	log.Printf("Running server on %s", port)

	http.Handle("/assets/", serveStaticFile("assets"))

	http.HandleFunc("/", index)

	http.ListenAndServe(port, nil)
}
