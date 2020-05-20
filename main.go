package main

import (
	"os"

	"github.com/urfave/cli"
	"github.com/verchol/kubectx/pkg/actions"
)

//AppVersion ...
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
	builtBy = "unknown"
)

func init() {
	app.Name = "Create Kubeconfig CLI"
	app.Usage = "An example how to create kube config"
	app.Author = "verchol"
	app.Version = version
}

var app = cli.NewApp()

func init() {

	actions.Commands(app)
}
func main() {

	//actions.commands(app)
	err := app.Run(os.Args)

	if err != nil {
		os.Exit(1)
	}
	return

}
