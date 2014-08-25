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
			Name:  "new",
			Usage: "Create new reka skeleton app",
			Action: func(c *cli.Context) {
				commands.Init(c.Args())
			},
		},
		{
			Name:  "server",
			Usage: "Add and remove servers",
			Subcommands: []cli.Command{
				{
					Name:  "add",
					Usage: "Register a new remote server",
					Action: func(c *cli.Context) {
						commands.ServerAdd(c.Args())
					},
				},
				{
					Name:  "ls",
					Usage: "List servers",
					Action: func(c *cli.Context) {
						commands.ServerList(c.Args())
					},
				},
			},
		},
	}

	app.Run(os.Args)

}
