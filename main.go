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

func index(w http.ResponseWriter, r *http.Request) {
	render("_views/index.tmpl", w)
}

func main() {
	var pastio Pastio

	var port string = fmt.Sprintf(":%s", pastio.Port())

	log.Printf("Running server on %s", port)

	http.HandleFunc("/", index)

	http.ListenAndServe(port, nil)
}
