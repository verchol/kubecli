package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"

	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"

	"io/ioutil"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"os"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

func init() {
	fmt.Println("in init")
	commands(app)
}
func homeDir() string {
	return "/Users/codefresh"
}

//SetNewCurrentContext ...
func SetNewCurrentContext(config clientcmd.ClientConfig, newContext string) clientcmd.ClientConfig {
	rawConfig, err := config.RawConfig()
	if err != nil {
		panic(err)
	}
	oldContext := rawConfig.CurrentContext
	for name := range rawConfig.Contexts {
		if name == newContext {
			rawConfig.CurrentContext = newContext
		}
	}

	if oldContext == rawConfig.CurrentContext {
		fmt.Printf("warning ,incorrect context %v\n", newContext)
		return config
	}
	fmt.Printf("%v\n", config.ConfigAccess())
	err = clientcmd.ModifyConfig(config.ConfigAccess(), rawConfig, false)
	if err != nil {
		panic(err)
	}
	newConfig := clientcmd.NewDefaultClientConfig(rawConfig, &clientcmd.ConfigOverrides{})
	return newConfig
}
func printAuth(auth *api.AuthInfo) {
	if auth.Token != "" {
		color.Green("[token] ")
		fmt.Printf("token = %v\n", auth.Token)
	}
	if auth.TokenFile != "" {
		color.Green("[tokenfile] ")
		fmt.Printf("%v\n", auth.TokenFile)
	}
	if auth.ClientCertificate != "" {
		color.Green("[cert] ")
		fmt.Printf("%v\n", auth.ClientCertificate)
		color.Green("[cert] ")
		fmt.Printf("%v\n", auth.ClientCertificateData)
	}
	if auth.AuthProvider != nil {
		fmt.Printf("[authProvider] = %v\n", auth.AuthProvider)
	}
}

func load() (*clientcmdapi.Config, error) {
	data, err := ioutil.ReadFile("./kubeconfig")
	if err != nil {
		panic(err)
	}
	config, err := clientcmd.Load(data)
	fmt.Printf("%v", config)

	return config, err
}

var app = cli.NewApp()

func info() {
	app.Name = "Create Config CLI"
	app.Usage = "An example how to create kube config"
	app.Author = "verchol"
	app.Version = "1.0.0"
}

//SetContextAction ....
func SetContextAction(c *cli.Context) error {

	fmt.Printf("command %v \n", c.Command)
	fmt.Printf("flags %v \n", c.App.Flags)
	fmt.Printf("args %v \n", c.Args())
	var newContext string

	config, err := loadConfig()
	if err != nil {
		return err
	}
	newContext = c.Args().First()
	if newContext != "" {
		fmt.Println("updating context...", newContext)
		config = SetNewCurrentContext(config, newContext)
	}

	green := color.New(color.FgGreen).SprintFunc()
	fmt.Printf("switched to context %v\n", green(newContext))

	return err
}
func handleSetContext(c *cli.Context) error {

	fmt.Printf("command %v \n", c.Command)
	fmt.Printf("flags %v \n", c.App.Flags)
	fmt.Printf("args %v \n", c.Args())
	var newContext string

	config, err := loadConfig()
	if err != nil {
		return err
	}
	if c.Command.Name == "newcontext" {
		newContext = c.String("context")
	}
	if newContext != "" {
		fmt.Println("updating context...", newContext)
		config = SetNewCurrentContext(config, newContext)
	}
	rawConfig, err := config.RawConfig()
	fmt.Printf("\ncurrent context is %v\n", rawConfig.CurrentContext)
	flags := FlagOptions{Validate: c.Bool("validate"), NoCache: c.Bool("nocache")}
	listContexts(config, flags)

	return nil
}
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
			Action: handleSetContext,
		},
		{
			Name:    "newcontext",
			Aliases: []string{"set"},
			Usage:   "set new context",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "context",
					Usage: "used to set new kubernetes context",
				}, cli.BoolFlag{
					Name:  "validate",
					Usage: "used to validate cluster connectivity",
				},
			},
			Action: handleSetContext,
		},
	}
	/*
		Name        string
		Usage       string
		EnvVar      string
		FilePath    string
		Required    bool
		Hidden      bool
		Destination *bool
	*/
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "validate",
			Usage: "used to validate cluster connectivity",
		},
	}
	app.Before = func(c *cli.Context) error {

		return nil
	}
	app.Action = SetContextAction
}
func loadConfig(opts ...string) (clientcmd.ClientConfig, error) {

	home, _ := os.UserHomeDir()
	var kubeconfig = filepath.Join(home, ".kube", "config")
	tempConfig, err := clientcmd.LoadFromFile(kubeconfig)
	green := color.New(color.FgGreen).SprintFunc()
	//tempConfig.CurrentContext
	config := clientcmd.NewDefaultClientConfig(*tempConfig, &clientcmd.ConfigOverrides{})
	fmt.Printf("%v\n", green(config.ConfigAccess().GetDefaultFilename()))

	return config, err

}
func listContexts(config clientcmd.ClientConfig, flags FlagOptions) {
	rawConfig, err := config.RawConfig()
	if err != nil {
		panic(err)
	}
	green := color.New(color.FgGreen).SprintFunc()

	var data [][]string
	headers := []string{"Contexts", "IsAvailable"}
	var contexts [][]string

	if !flags.Validate {
		headers = headers[:1]
	}
	//load cache

	bytes, err := ioutil.ReadFile(".status-cache")
	if err != nil {
		fmt.Println(green("no status file yet created\n"))

	}
	if err == nil && !flags.NoCache {
		err := json.Unmarshal(bytes, &contexts)
		if err != nil {
			panic(err)
		}
		fmt.Println("context from cache\n")
		printTable(headers, contexts)
		return
	}

	for name := range rawConfig.Contexts {
		var pods int
		pods = -1
		var contextData []string
		currentContextName := fmt.Sprintf("%v", name)

		if rawConfig.CurrentContext == name {
			//fmt.Printf("current context %v is active =  %v\n", green(name), (pods != -1))
			currentContextName = fmt.Sprintf("%v", green(name))

		}

		var podStr string
		if flags.Validate {
			pods = testCluster(config, name)
		}
		if pods != -1 {
			podStr = fmt.Sprintf("%v", pods)
		} else {
			podStr = "_"
		}
		contextData = []string{currentContextName, podStr}
		data = append(data, contextData)

	}

	printTable(headers, data)
	if flags.Validate {
		bytes, err = json.Marshal(data)

		if err != nil {
			panic(err)
		}
		err = ioutil.WriteFile(".status-cache", bytes, 0644)
		if err != nil {
			panic(err)
		}
	}

}

//Contexts
type Contexts [][]string

func (a Contexts) Len() int { return len(a) }
func (a Contexts) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a Contexts) Less(i, j int) bool {

	return a[i][0] < a[j][0]
}

func printTable(header []string, data [][]string) {

	sort.Sort(Contexts(data))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)

	for _, v := range data {
		table.Append(v)
	}
	table.Render() // Send output
}
func setNewDefaultContext(ctx string) {}
func testCluster(config clientcmd.ClientConfig, currentContext string) int {
	r, _ := config.RawConfig()
	tempConfig := clientcmd.NewDefaultClientConfig(r,
		&clientcmd.ConfigOverrides{CurrentContext: currentContext})

	restConfig, err := tempConfig.ClientConfig()
	if err != nil {
		return -1
	}
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return -1
	}
	ns, _, _ := config.Namespace()
	var timeout int64 = 5

	fmt.Printf("checking if cluster %v avaialable\n", restConfig.Host)

	pods, err :=
		clientset.
			CoreV1().
			Pods(ns).
			List(metav1.ListOptions{TimeoutSeconds: &timeout})

	fmt.Printf("[%v] pods are %v len=%v \n", currentContext, pods.Items, len(pods.Items))
	if err != nil {
		return -1
	}

	return len(pods.Items)

}

type FlagOptions struct {
	List     bool
	Validate bool
	NoCache  bool
}

func main() {

	commands(app)
	/*
		command := []string{"newcontext", "--context"}
		first := make([]string, 1)

		copy(first, os.Args[:1])
		first = append(first, command...)
		os.Args = append(first, os.Args[1:]...)
		fmt.Printf("args %v", os.Args)*/

	err := app.Run(os.Args)

	if err != nil {
		panic(err)
	}
	return

}
