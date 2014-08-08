package confutil

import (
	"io/ioutil"
	"log"
	"os"
)

func WriteIfMissing(path string, content string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		Write(path, content)
	}
}

func Write(path string, content string) {
	err := ioutil.WriteFile(path, []byte(content), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
