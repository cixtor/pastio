package main

import (
	"flag"
	"log"
	"time"

	"github.com/cixtor/middleware"
)

// PublicFolder is the document root.
const PublicFolder = "assets"

// StorageFolder is where the submissions will be stored.
const StorageFolder = "storage"

var app Application

var router = middleware.New()

var addr string

func main() {
	flag.StringVar(&addr, "addr", ":3000", "Hostname and port number to listen for HTTP requests")
	flag.Parse()

	router.ReadTimeout = time.Second * 5
	router.WriteTimeout = time.Second * 10

	router.STATIC(PublicFolder, "/assets")

	router.POST("/save", app.Save)
	router.GET("/modes", app.Modes)
	router.GET("/raw/:unique", app.RawCode)
	router.GET("/", app.Index)

	if err := router.ListenAndServe(addr); err != nil {
		log.Fatal(err)
	}
}
