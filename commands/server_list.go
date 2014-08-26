package commands

import (
	"github.com/wsxiaoys/terminal/color"
	"log"
	"reka/config"
	"strconv"
)

func ServerList(Args []string) {
	config, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	length := 0
	for name := range config.Servers {
		if len(name) > length {
			length = len(name)
		}
	}
	for name, server := range config.Servers {
		if name == config.DefaultServer {
			color.Printf((" * @{!}%" + strconv.Itoa(length) + "s@{|} %s\n"), name, server.URL)
		} else {
			color.Printf(("   @{!}%" + strconv.Itoa(length) + "s@{|} %s\n"), name, server.URL)
		}
	}
}
