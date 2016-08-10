package main

import (
	"net/http"
)

func index(w http.ResponseWriter, r *http.Request) {
	render("_views/index.tmpl", w)
}
