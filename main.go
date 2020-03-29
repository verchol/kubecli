package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/urfave/cli"
	"github.com/verchol/kubectx/pkg/config"
)

func info() {
	app.Name = "Create Kubeconfig CLI"
	app.Usage = "An example how to create kube config"
	app.Author = "verchol"
	app.Version = "1.0.0"
}

var app = cli.NewApp()

func commands(app *cli.App) {
	app.Commands = []cli.Command{
		{
			Name:    "version",
			Aliases: []string{"version"},
			Usage:   "create a new config",
			Action: func(c *cli.Context) {
				color.Green("cli info is %v\n", app.Usage)
			},
		},
		{
			Name:   "delete",
			Usage:  "delete context",
			Action: config.HandleDeleteContext,
		},
		{
			Name:   "switch",
			Usage:  "change context to new one",
			Action: config.HandleSetContext,
		},
		{
			Name:  "list",
			Usage: "list context",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "validate",
					Usage: "used to validate cluster connectivity",
				},
				cli.BoolFlag{
					Name:  "nocache",
					Usage: "used to reinitiate conncectivity status",
				},
			},
			Action: config.HandleSetContext,
		},
		{
			Name:    "newcontext",
			Aliases: []string{"new"},
			Usage:   "set new context",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Usage: "used to set new kubernetes context",
				}, cli.BoolFlag{
					Name:     "token",
					Usage:    "used to define context's service account ",
					Required: false,
				},
				cli.StringFlag{
					Name:     "namespace",
					Usage:    "used to define context's namespace",
					Required: true,
				},
			},
			Action: config.CreateContextAction,
		},
	}

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "validate",
			Usage: "used to validate cluster connectivity",
		}, cli.BoolFlag{
			Name:  "nocache",
			Usage: "used to reinitiate conncectivity status",
		},
	}

	app.Action = config.SetContextAction
}

func init() {

	commands(app)
}
func main() {

	commands(app)
	err := app.Run(os.Args)

	if err != nil {
		os.Exit(1)
	}
	return

}
