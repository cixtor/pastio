package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/cixtor/middleware"
)

// Application serves as the base for all the API endpoints
// exposed by the web server. Some of the associated methods
// write the content of a dynamic template and others are used
// to process a HTTP request and return either a valid JSON
// object or a HTTP status code.
type Application struct{}

// Response is the basic JSON object returned by the API.
type Response struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	Metadata string `json:"metadata"`
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

// Save accepts a POST request with the mode and code
//
// The parameters mode and code contain the syntax highlighter used by the user
// during the creation of the text that will be sent to the server and the code
// in itself respectively. Notice that additional information like the origin IP
// address, request timestamp, visibility of the code, and referer will be
// recorded along with the syntax and the content.
//
// The visibility is intended to be used by the code viewer to determine if the
// IP has access to the code, for example, when someone submits a code and
// selects "private" the code will only be accessible if the IP of the request
// is in the allowed list.
func (app *Application) Save(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Invalid HTTP method", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		log.Println(err)
		return
	}

	var content string
	mode := r.Form.Get("mode")
	code := r.Form.Get("code")

	if code == "" {
		http.Error(w, "Cannot save empty text", http.StatusBadRequest)
		return
	}

	if !modeExists(mode) {
		http.Error(w, "Syntax does not exists", http.StatusBadRequest)
		return
	}

	content += "=== start_metadata\n"
	content += fmt.Sprintf("remote_addr: %s\n", r.RemoteAddr)
	content += fmt.Sprintf("request_time: %d\n", int32(time.Now().Unix()))
	content += fmt.Sprintf("referer: %s\n", r.Header.Get("Referer"))
	content += fmt.Sprintf("visibility: %s\n", r.Form.Get("visibility"))
	content += fmt.Sprintf("mode: %s\n", mode)
	content += "=== end_metadata\n"
	content += code

	fpath, fname := uniqueFname(6)

	if err := saveFile(fpath, content); err != nil {
		log.Println("SaveFile:", err)
		http.Error(w, "Could not save content", http.StatusInternalServerError)
		return
	}

	var success Response

	success.Status = "ok"
	success.Message = "Operation was successful"
	success.Metadata = string(fname)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewEncoder(w).Encode(success); err != nil {
		log.Println(err)
		return
	}
}
