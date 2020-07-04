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
	app.Usage = "kubernetes context management utility"
	app.Author = "Oleg Verhovsky"
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
		panic(err)
	}
	return

}
