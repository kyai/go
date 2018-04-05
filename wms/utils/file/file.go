package file

import (
	"io/ioutil"
	"os"
)

func Reader(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	text, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(text), err
}
