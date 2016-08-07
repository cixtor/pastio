package main

import (
	"fmt"
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

func Index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello world!\n"))
}

func main() {
	var pastio Pastio

	var port string = fmt.Sprintf(":%s", pastio.Port())

	log.Printf("Running server on %s", port)

	http.HandleFunc("/", Index)

	http.ListenAndServe(port, nil)
}
