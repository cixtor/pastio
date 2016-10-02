package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type Application struct{}

type Response struct {
	Status   string `json:"status"`
	Message  string `json:"message"`
	Metadata string `json:"metadata"`
}

type ModeList struct {
	Status  string   `json:"status"`
	Default string   `json:"default"`
	Modes   []string `json:"modes"`
}

func (app *Application) Index(w http.ResponseWriter, r *http.Request) {
	render("_views/index.tmpl", w)
}

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

	file, err := os.Open(fpath)
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		w.Write([]byte(scanner.Text() + "\n"))
	}
}

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
