package commands

import (
	"github.com/nicksellen/reka/config"
	"github.com/wsxiaoys/terminal/color"
	"log"
	"strconv"
)

func DeploymentList(Args []string) {
	config, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	length := 0
	for name := range config.Deployments {
		if len(name) > length {
			length = len(name)
		}
	}
	for name, deployment := range config.Deployments {
		if name == config.DefaultDeployment {
			color.Printf((" * @{!}%" + strconv.Itoa(length) + "s@{|} %s\n"), name, deployment.Url)
		} else {
			color.Printf(("   @{!}%" + strconv.Itoa(length) + "s@{|} %s\n"), name, deployment.Url)
		}
	}
}
