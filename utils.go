package main

import (
	"io/ioutil"
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
