package main

import (
	"encoding/json"
	"log"
	"net/http"
)

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

	success := struct {
		Status  string   `json:"status"`
		Default string   `json:"default"`
		Modes   []string `json:"modes"`
	}{
		Status:  "ok",
		Default: "php",
		Modes:   availableModes,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err := json.NewEncoder(w).Encode(success); err != nil {
		log.Println(err)
		return
	}
}

func modeExists(mode string) bool {
	for _, name := range availableModes {
		if name == mode {
			return true
		}
	}
	return false
}

var availableModes = []string{
	"abap",
	"actionscript",
	"ada",
	"asciidoc",
	"assembly_x86",
	"autohotkey",
	"batchfile",
	"c9search",
	"c_cpp",
	"clojure",
	"cobol",
	"coffee",
	"coldfusion",
	"csharp",
	"css",
	"curly",
	"dart",
	"diff",
	"django",
	"d",
	"dot",
	"ejs",
	"erlang",
	"forth",
	"ftl",
	"glsl",
	"golang",
	"groovy",
	"haml",
	"handlebars",
	"haskell",
	"haxe",
	"html_completions",
	"html",
	"html_ruby",
	"ini",
	"jack",
	"jade",
	"java",
	"javascript",
	"jsoniq",
	"json",
	"jsp",
	"jsx",
	"julia",
	"latex",
	"less",
	"liquid",
	"lisp",
	"livescript",
	"logiql",
	"lsl",
	"lua",
	"luapage",
	"lucene",
	"makefile",
	"markdown",
	"matlab",
	"mushcode_high_rules",
	"mushcode",
	"mysql",
	"nix",
	"objectivec",
	"ocaml",
	"pascal",
	"perl",
	"pgsql",
	"php",
	"plain_text",
	"powershell",
	"prolog",
	"properties",
	"protobuf",
	"python",
	"rdoc",
	"rhtml",
	"r",
	"ruby",
	"rust",
	"sass",
	"scad",
	"scala",
	"scheme",
	"scss",
	"sh",
	"sjs",
	"snippets",
	"soy_template",
	"space",
	"sql",
	"stylus",
	"svg",
	"tcl",
	"tex",
	"textile",
	"text",
	"toml",
	"twig",
	"typescript",
	"vbscript",
	"velocity",
	"verilog",
	"xml",
	"xquery",
	"yaml",
}
