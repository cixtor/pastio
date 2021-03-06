package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
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

// ModeList is the JSON object for the supported editor syntax.
type ModeList struct {
	Status  string   `json:"status"`
	Default string   `json:"default"`
	Modes   []string `json:"modes"`
}

// Index parses and renders the template for the homepage.
func (app *Application) Index(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("_views/index.tmpl")

	if err != nil {
		log.Println("TPL: index.tmpl;", err)
		return
	}

	tpl.Execute(w, nil)
}

// Modes is the API endpoint resposible for returning ModeList
// which is intended to contain a list of syntax highlighters of
// popular programming languages. Every time the user clicks one
// of the options the code editor will change its syntax
// highlighter to increase the readability of the code.
func (app *Application) Modes(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Invalid HTTP method", http.StatusBadRequest)
		return
	}

	var success ModeList

	success.Status = "ok"
	success.Default = "php"
	success.Modes = availableModes()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	json.NewEncoder(w).Encode(success)
}

// RawCode searches and renders the content of the requested
// file in plain/text. The content of the file will not be
// altered, it will be shown as is. If the unique ID is not
// found in the storage folder the server will respond with a
// 404 Not Found status code. Notice that the metadata will be
// omitted from the response.
func (app *Application) RawCode(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()
	var unique = ctx.Value("unique")

	if unique == nil {
		http.Error(w, "Unique ID is missing", http.StatusBadRequest)
		return
	}

	fpath, err := fullFpath(unique.(string))

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")

	var line string
	var safeToPrint bool
	file, err := os.Open(fpath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = scanner.Text()
		if line == "=== end_metadata" {
			safeToPrint = true
			continue
		}
		if safeToPrint {
			w.Write([]byte(line + "\n"))
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

	r.ParseForm()

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

	json.NewEncoder(w).Encode(success)
}
