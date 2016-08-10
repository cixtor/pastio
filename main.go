package main

import (
	"log"
	"net/http"
	"os"
)

const DEFAULT_PORT = "8080"
const PUBLIC_FOLDER = "assets"

func main() {
	var port string = os.Getenv("PORT")

	/**
	 * Use default port number
	 *
	 * If no custom port number is specified via an environment variable we will
	 * use the default one, hoping that it is not being used by a difference
	 * process, the program will panic otherwise.
	 *
	 * @type {string}
	 */
	if port == "" {
		port = DEFAULT_PORT
	}

	log.Printf("Running server on :%s", port)

	http.Handle("/assets/", serveStaticFile(PUBLIC_FOLDER))

	http.HandleFunc("/", index)

	http.ListenAndServe(":"+port, logger(http.DefaultServeMux))
}
