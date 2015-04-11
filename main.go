package main

import (
	"github.com/codegangsta/cli"
	"github.com/nicksellen/reka/commands"
	"os"
)

func main() {

	app := cli.NewApp()

	app.Name = "reka"
	app.Usage = "cli for using reka :)"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:  "deploy",
			Usage: "Deploy app to a server",
			Action: func(c *cli.Context) {
				commands.Deploy(c.Args())
			},
		},
		{
			Name:  "undeploy",
			Usage: "Undeploy app from a server",
			Action: func(c *cli.Context) {
				commands.Undeploy(c.Args())
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
			Name:  "deployment",
			Usage: "add/remove/list deployments",
			Subcommands: []cli.Command{
				{
					Name:  "add",
					Usage: "register a deployment",
					Action: func(c *cli.Context) {
						commands.DeploymentAdd(c.Args())
					},
				},
				{
					Name:  "list",
					Usage: "list deployments",
					Action: func(c *cli.Context) {
						commands.DeploymentList(c.Args())
					},
				},
			},
		},
	}

	app.Run(os.Args)

}
