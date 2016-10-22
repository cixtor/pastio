package main

import "os"
import "github.com/cixtor/middleware"

// PublicFolder is the document root.
const PublicFolder = "assets"

// StorageFolder is where the submissions will be stored.
const StorageFolder = "storage"

func main() {
	var app Application

	router := middleware.New()

	router.ReadTimeout = 5
	router.WriteTimeout = 10
	router.Port = os.Getenv("PORT")

	router.STATIC(PublicFolder, "/assets")

	router.POST("/save", app.Save)
	router.GET("/modes", app.Modes)
	router.GET("/raw/:unique", app.RawCode)
	router.GET("/", app.Index)

	router.ListenAndServe()
}
