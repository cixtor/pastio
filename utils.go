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
	result := make([]byte, length)
	alpha := []byte("abcdefghijklmnopqrstuvwxyz")
	total := len(alpha)

	for i := 0; i < length; i++ {
		result[i] = alpha[rand.Intn(total)]
	}

	fpath := pwd(StorageFolder + "/" + string(result[0]) + "/" + string(result) + ".txt")

	if fileExists(fpath) {
		return uniqueFname(length)
	}

	return fpath, result
}
