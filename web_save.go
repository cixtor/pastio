package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cixtor/middleware"
)

func init() {
	router.POST("/save", app.Save)
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
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		return
	}

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

	content := "=== start_metadata\n"
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

	success := struct {
		Status   string `json:"status"`
		Message  string `json:"message"`
		Metadata string `json:"metadata"`
	}{
		Status:   "ok",
		Message:  "Operation was successful",
		Metadata: string(fname),
	}

	if err := middleware.JSON(w, r, success); err != nil {
		log.Println(err)
	}
}
