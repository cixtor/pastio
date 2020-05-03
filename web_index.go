package main

import (
	"html/template"
	"log"
	"net/http"
)

func init() {
	router.GET("/", app.Index)
}

// Index parses and renders the template for the homepage.
func (app *Application) Index(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("_views/index.tmpl")

	if err != nil {
		log.Println("TPL: index.tmpl;", err)
		return
	}

	if err := tpl.Execute(w, nil); err != nil {
		log.Println(err)
		return
	}
}
