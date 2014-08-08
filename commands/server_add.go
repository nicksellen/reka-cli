package commands

import (
	"log"
	"path"
	"reka/config/confutil"
)

func ServerAdd(Args []string) {
	if len(Args) != 2 {
		log.Fatal("please provide server name and address")
	}
	name := Args[0]
	url := Args[1]
	confutil.Write(path.Join(".reka/config/servers", name), url)
	confutil.WriteIfMissing(".reka/config/default-server", name)
}
