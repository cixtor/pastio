package main

import (
	"bufio"
	"net/http"
	"os"

	"github.com/cixtor/middleware"
)

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

	var line string
	var safeToPrint bool

	file, err := os.Open(fpath)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	defer file.Close()

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
