package main

import "os"
import "github.com/cixtor/middleware"

const PUBLIC_FOLDER = "assets"
const STORAGE_FOLDER = "storage"

func main() {
	var app Application

	router := middleware.New()

	router.ReadTimeout = 5
	router.WriteTimeout = 10
	router.Port = os.Getenv("PORT")

	router.STATIC(PUBLIC_FOLDER, "/assets")

	router.POST("/save", app.Save)
	router.GET("/modes", app.Modes)
	router.GET("/raw/:unique", app.RawCode)
	router.GET("/", app.Index)

	router.ListenAndServe()
}
