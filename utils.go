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
	directory := filename[0 : len(filename)-10]

	if err := os.MkdirAll(directory, 0755); err != nil {
		return err
	}

	return ioutil.WriteFile(filename, []byte(data), 0644)
}

func fileExists(fpath string) bool {
	_, err := os.Stat(fpath)

	return err == nil
}

func uniqueFname(length int) (string, []byte) {
	var total int
	var fpath string
	var result = make([]byte, length)
	var alpha = []byte("abcdefghijklmnopqrstuvwxyz")

	total = len(alpha)

	for i := 0; i < length; i++ {
		result[i] = alpha[rand.Intn(total)]
	}

	fpath += StorageFolder + "/"
	fpath += string(result[0]) + "/"
	fpath += string(result) + ".txt"

	fpath = pwd(fpath)

	if fileExists(fpath) {
		return uniqueFname(length)
	}

	return fpath, result
}
