package main

import (
	"errors"
	"io/ioutil"
	"math/rand"
	"os"
)

func pwd(trailing string) string {
	dir, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	return dir + "/" + trailing
}

func saveFile(filename string, data string) error {
	var directory string = filename[0 : len(filename)-10]

	if err := os.MkdirAll(directory, 0755); err != nil {
		return err
	}

	return ioutil.WriteFile(filename, []byte(data), 0644)
}

func fileExists(fpath string) bool {
	_, err := os.Stat(fpath)

	return err == nil
}

func modeExists(mode string) bool {
	return fileExists(pwd(PUBLIC_FOLDER + "/js/ace/mode-" + mode + ".js"))
}

func fullFpath(unique string) (string, error) {
	var fpath string

	fpath += STORAGE_FOLDER + "/"
	fpath += string(unique[0]) + "/"
	fpath += string(unique) + ".txt"

	if !fileExists(fpath) {
		return "", errors.New("File does not exists")
	}

	return fpath, nil
}

func uniqueFname(length int) (string, []byte) {
	var fpath string
	var result = make([]byte, length)
	var alpha = []byte("abcdefghijklmnopqrstuvwxyz")
	var total int = len(alpha)

	for i := 0; i < length; i++ {
		result[i] = alpha[rand.Intn(total)]
	}

	fpath += STORAGE_FOLDER + "/"
	fpath += string(result[0]) + "/"
	fpath += string(result) + ".txt"

	fpath = pwd(fpath)

	if fileExists(fpath) {
		return uniqueFname(length)
	}

	return fpath, result
}

func availableModes() []string {
	return []string{
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
}
