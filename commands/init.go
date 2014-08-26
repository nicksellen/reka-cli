package commands

import (
	"log"
	"os"
)

func Init(Args []string) {

	if len(Args) != 1 {
		log.Fatal("please provide the directory to initialize")
	}

	dir := Args[0]

	info, err := os.Stat(dir)

	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}

	if dir != "." && info != nil && info.IsDir() {
		log.Fatal(dir + " already exists")
	}

	os.MkdirAll(dir+"/.reka/config/servers", 0755)

	main := dir + "/main.reka"

	_, err = os.Stat(main)
	if err != nil {
		if os.IsNotExist(err) {
			file, err := os.Create(main)
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()
		} else {
			log.Fatal(err)
		}
	}

}
