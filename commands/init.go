package commands

import (
	"log"
	"os"
)

func Init(Args []string) {

	if len(Args) != 1 {
		log.Fatal("please provide the application name")
	}

	identity := Args[0]

	info, err := os.Stat(identity)
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}

	if info != nil && info.IsDir() {
		log.Fatal(identity + " already exists")
	}

	os.MkdirAll(identity+"/.reka/config/servers", 0755)
	file, err := os.Create(identity + "/main.reka")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

}
