package main

import (
	"github.com/codegangsta/cli"
	"os"
	"reka/commands"
  _ "crypto/sha512"
)

func main() {

	app := cli.NewApp()

	app.Name = "reka"
	app.Usage = "cli for using reka :)"
	app.Version = "0.0.1"

	app.Commands = []cli.Command{
		{
			Name:  "send",
			Usage: "Upload files to reka, ready for deployment",
			Action: func(c *cli.Context) {
				commands.Upload(c.Args())
			},
		},
		{
			Name:  "validate",
			Usage: "Validate some previously upload config",
			Action: func(c *cli.Context) {
				commands.Validate(c.Args())
			},
		},
		{
			Name:  "deploy",
			Usage: "Deploy some previously upload config",
			Action: func(c *cli.Context) {
				commands.Deploy(c.Args())
			},
		},
		{
			Name:  "redeploy",
			Usage: "Redeploy an application by uuid",
			Action: func(c *cli.Context) {
				commands.Redeploy(c.Args())
			},
		},
		{
			Name:  "undeploy",
			Usage: "Undeploy an application by uuid",
			Action: func(c *cli.Context) {
				commands.Undeploy(c.Args())
			},
		},
	}

	app.Run(os.Args)

}
