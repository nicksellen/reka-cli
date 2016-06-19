package commands

import (
	"github.com/nicksellen/reka-cli/config/confutil"
	"log"
	"path"
)

func DeploymentAdd(Args []string) {
	if len(Args) != 2 {
		log.Fatal("please provide name and deployment address")
	}
	name := Args[0]
	url := Args[1]
	confutil.Write(path.Join(".reka/config/deployments", name), url)
	confutil.WriteIfMissing(".reka/config/default-deployment", name)
}
