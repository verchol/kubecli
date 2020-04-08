package main

import (
	"os"

	"github.com/urfave/cli"
	"github.com/verchol/kubectx/pkg/actions"
)

func info() {
	app.Name = "Create Kubeconfig CLI"
	app.Usage = "An example how to create kube config"
	app.Author = "verchol"
	app.Version = "1.0.0"
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
