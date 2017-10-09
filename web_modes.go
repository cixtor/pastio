package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// ModeList is the JSON object for the supported editor syntax.
type ModeList struct {
	Status  string   `json:"status"`
	Default string   `json:"default"`
	Modes   []string `json:"modes"`
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

	if err := json.NewEncoder(w).Encode(success); err != nil {
		log.Println(err)
		return
	}
}
