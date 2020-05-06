package main

import (
	"bufio"
	"errors"
	"net/http"
	"os"

	"github.com/cixtor/middleware"
)

func init() {
	router.GET("/raw/:unique", app.RawCode)
}

// RawCode searches and renders the content of the requested
// file in plain/text. The content of the file will not be
// altered, it will be shown as is. If the unique ID is not
// found in the storage folder the server will respond with a
// 404 Not Found status code. Notice that the metadata will be
// omitted from the response.
func (app *Application) RawCode(w http.ResponseWriter, r *http.Request) {
	unique := middleware.Param(r, "unique")

	fpath, err := fullFpath(unique)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

	file, err := os.Open(fpath)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	defer file.Close()

	var line string
	var safeToPrint bool

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line = scanner.Text()
		if line == "=== end_metadata" {
			safeToPrint = true
			continue
		}
		if safeToPrint {
			_, _ = w.Write([]byte(line + "\n"))
		}
	}
}

func fullFpath(unique string) (string, error) {
	fpath := storageFolder + "/" + string(unique[0]) + "/" + unique + ".txt"

	if !fileExists(fpath) {
		return "", errors.New("File does not exists")
	}

	return fpath, nil
}
