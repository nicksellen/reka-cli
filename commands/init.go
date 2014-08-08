package commands

import (
	"log"
	"os"
	"reka/config/confutil"
)

func Init(Args []string) {

	if len(Args) != 1 {
		log.Fatal("please provide the application identity")
	}

	identity := Args[0]

	os.MkdirAll(".reka/config/servers", 0755)

	confutil.WriteIfMissing(".reka/config/identity", identity)
	confutil.WriteIfMissing(".reka/config/work", ".")

}
