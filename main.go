package main

import (
	"github.com/codegangsta/cli"
	"os"
	"reka/commands"
)

func main() {

	app := cli.NewApp()

	app.Name = "reka"
	app.Usage = "cli for using reka :)"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:  "push",
			Usage: "Push to reka instance to run",
			Action: func(c *cli.Context) {
				commands.Push(c.Args())
			},
		},
		{
			Name:  "init",
			Usage: "Initialize reka config in current dir",
			Action: func(c *cli.Context) {
				commands.Init(c.Args())
			},
		},
		{
			Name:  "server-add",
			Usage: "Add a server",
			Action: func(c *cli.Context) {
				commands.ServerAdd(c.Args())
			},
		},
	}

	app.Run(os.Args)

}
