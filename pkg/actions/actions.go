package actions

import (
	"fmt"

	"github.com/docker/machine/libmachine/log"
	"github.com/urfave/cli"
	"github.com/verchol/kubectx/pkg/config"
	configtools "github.com/verchol/kubectx/pkg/config"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/fatih/color"
)

//Commands ...
func Commands(app *cli.App) {
	app.Commands = []cli.Command{
		{
			Name:    "version",
			Aliases: []string{"version"},
			Usage:   "create a new config",
			Action: func(c *cli.Context) {
				color.Green("cli info is %v\n", app.Version)
			},
		},
		{
			Name:    "test",
			Aliases: []string{"t"},
			Usage:   "test cluster",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "context",
					Usage: "context to for validation",
				},
				cli.Int64Flag{
					Name:  "timeout",
					Usage: "how long to wait for cluster to answer in sec",
				},
			},
			Action: TestClusterAction,
		},
		{
			Name:    "namespace",
			Aliases: []string{"ns"},
			Usage:   "set new namespace for context",
			Action:  SetContextNamespaceAction,
		},
		{
			Name:   "delete",
			Usage:  "delete context",
			Action: configtools.HandleDeleteContext,
		},
		{
			Name:    "switch",
			Aliases: []string{"s"},
			Usage:   "change context to new one",
			Action:  configtools.HandleSetContext,
		},
		{
			Name:    "list",
			Usage:   "list context",
			Aliases: []string{"ls"},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "validate",
					Usage: "used to validate cluster connectivity",
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
				cli.StringSliceFlag{
					Name:     "verbs",
					Usage:    "used to define what role to use",
					Required: false,
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
			Name:  "cache",
			Usage: "take information from local cache file ",
		},
	}

	app.Action = configtools.SetContextAction
}

//CreqteContextAction ...
func CreateContextAction(c *cli.Context) error {

	contextName := c.String("name")
	if contextName == "" {
		contextName = c.Args().First()
	}
	ns := c.String("namespace")
	sa := fmt.Sprintf("sa-%v-%v", ns, contextName)
	//_ := c.Bool("create")
	//_ := c.String("permission")

	log.Info(contextName, ns)
	log.Info("verbs %v\n", c.StringSlice("verbs"))

	config, err := configtools.LoadConfig()
	if err != nil {
		return err
	}
	roleOpts := configtools.NewRoleOpts(fmt.Sprintf("role-%v-%v", sa, ns), ns)
	role, err := configtools.CreateRole(roleOpts, config)

	if err != nil {

		//TODO identify when the error reason is "AlreadyExists"
		//For now is skipping treating error

		log.Error(err)
	}

	saObj, err := configtools.CreateServiceAccount(ns, sa, config)

	if err != nil {
		//TODO identify when the error reason is "AlreadyExists"
		//For now is skipping treating error

		log.Error(err)
	}

	roleBindingOpts := configtools.NewRoleBindingOpts(fmt.Sprintf("rb1-%v-%v", sa, ns), ns)
	roleBindingOpts.Role = role.Name
	roleBindingOpts.ServiceAccount = sa
	roleBindingOpts.ServiceAccountNs = ns

	_, err = configtools.CreateRoleBinding(roleBindingOpts, config)
	if err != nil {
		//TODO identify when the error reason is "AlreadyExists"
		//For now is skipping treating error

		log.Error(err)
	}

	err = configtools.CreateContext(contextName, ns, string(saObj.Token), config)

	return err

}

//TestClusterAction ...
func TestClusterAction(c *cli.Context) error {
	context := c.String("context")

	if context == "" {
		context = c.Args().First()
	}

	fmt.Printf("context is %v\n", context)

	config, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	r, _ := config.RawConfig()
	if context == "" {
		context = r.CurrentContext
	}
	waitingPeriod := c.Int64("timeout")

	tempConfig := clientcmd.NewDefaultClientConfig(r,
		&clientcmd.ConfigOverrides{CurrentContext: context})

	namespace, _, err := tempConfig.Namespace()
	if err != nil {
		panic(err)
	}

	restConfig, err := tempConfig.ClientConfig()

	if err != nil {
		panic(err)
	}
	clientSet, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		panic(err)
	}

	works, err := configtools.ValidateCluster(waitingPeriod, namespace, clientSet)
	Red := color.New(color.FgRed).SprintFunc()
	Green := color.New(color.FgGreen).SprintFunc()

	cache, err := configtools.NewLocalCache()

	if err != nil {
		panic(err)
	}
	kubeContext := configtools.KubeContext{}
	kubeContext.Name = context
	ns, _, _ := tempConfig.Namespace()
	kubeContext.Namespace = ns
	rawConfig, _ := config.RawConfig()
	auth := rawConfig.Contexts[context].AuthInfo
	authProvider := rawConfig.AuthInfos[auth].AuthProvider
	if authProvider != nil {
		kubeContext.AuthProvider = authProvider.Name
	}
	kubeContext.Status = configtools.ClusterAvailable

	if !works {
		fmt.Printf("context %v is not available \n%v :  %v\n", Green(context), Red("error:"), err)
		kubeContext.Status = configtools.ClusterNotAvailable

	} else {
		fmt.Printf("context %v is available\n", Green(context))
		kubeContext.Status = configtools.ClusterAvailable
	}

	fmt.Printf("adding entry %v %v\n", kubeContext.Name, kubeContext.Status)

	cache.AddEntry(context, &kubeContext)
	cache.Flash()
	return err
}

func SetContextNamespaceAction(c *cli.Context) error {

	ns := c.Args().First()

	config, err := config.LoadConfig()
	if err != nil {
		return err
	}

	rawConfig, err := config.RawConfig()
	if err != nil {
		return err
	}
	currentCtxName := rawConfig.CurrentContext
	context := rawConfig.Contexts[currentCtxName]
	context.Namespace = ns

	err = clientcmd.ModifyConfig(config.ConfigAccess(), rawConfig, false)

	fmt.Printf("context %v was updated with namespace %v \n", currentCtxName, ns)

	return err
}
