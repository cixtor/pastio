package main

import (
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
