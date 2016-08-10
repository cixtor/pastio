package main

import (
	"html/template"
	"log"
	"net/http"
)

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
