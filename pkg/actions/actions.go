package actions

import (
	"fmt"

	"github.com/docker/machine/libmachine/log"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	configtools "github.com/verchol/kubectx/pkg/config"
)

//Commands ...
func Commands(app *cli.App) {
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
			Action: configtools.HandleDeleteContext,
		},
		{
			Name:   "switch",
			Usage:  "change context to new one",
			Action: configtools.HandleSetContext,
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
			Action: configtools.HandleSetContext,
		},
		{
			Name:    "newcontext",
			Aliases: []string{"ctx"},
			Usage:   "set new context",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "name",
					Usage: "used to set new kubernetes context",
				}, cli.BoolFlag{
					Name:     "token",
					Usage:    "used to define context's service account ",
					Required: false,
				}, cli.StringFlag{
					Name:     "serviceAccount",
					Usage:    "usedservice account name",
					Required: false,
				},
				cli.StringFlag{
					Name:     "namespace",
					Usage:    "used to define context's namespace",
					Required: true,
				},
			},
			Action: CreateContextAction,
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

	app.Action = configtools.SetContextAction
}

//CreqteContextAction ...
func CreateContextAction(c *cli.Context) error {

	contextName := c.String("name")
	ns := c.String("namespace")
	sa := fmt.Sprintf("sa-%v-%v", ns, contextName)
	//_ := c.Bool("create")
	//_ := c.String("permission")

	log.Info(contextName, ns)
	config, err := configtools.LoadConfig()
	if err != nil {
		return err
	}
	roleOpts := configtools.NewRoleOpts(fmt.Sprintf("role-%v-%v", sa, ns), ns)
	role, err := configtools.CreateRole(roleOpts, config)

	if err != nil {
		panic(err)
	}

	saObj, err := configtools.CreateServiceAccount(ns, sa, config)

	if err != nil {
		return err
	}

	roleBindingOpts := configtools.NewRoleBindingOpts(fmt.Sprintf("rb1-%v-%v", sa, ns), ns)
	roleBindingOpts.Role = role.Name
	roleBindingOpts.ServiceAccount = sa
	roleBindingOpts.ServiceAccountNs = ns

	_, err = configtools.CreateRoleBinding(roleBindingOpts, config)
	if err != nil {
		panic(err)
	}

	err = configtools.CreateContext(contextName, ns, string(saObj.Token), config)

	return err

}
