package main

import (
	"flag"
	"log"
	"time"

	"github.com/cixtor/middleware"
)

var app Application

var router = middleware.New()

var addr string
var publicHTML string
var storageFolder string

func main() {
	flag.StringVar(&publicHTML, "public-html", "assets", "Directory with all the js, css, and image files")
	flag.StringVar(&storageFolder, "storage-folder", "storage", "Directory where the submissions will be stored")
	flag.StringVar(&addr, "addr", ":3000", "Hostname and port number to listen for HTTP requests")
	flag.Parse()

	router.ReadTimeout = time.Second * 5
	router.WriteTimeout = time.Second * 10

	router.STATIC(publicHTML, "/assets")

	if err := router.ListenAndServe(addr); err != nil {
		log.Fatal(err)
	}
}
