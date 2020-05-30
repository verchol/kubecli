package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"os"

	"github.com/docker/machine/libmachine/log"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
)

func HomeDir() string {
	return "/Users/codefresh"
}

//SetNewCurrentContext ...
func SetNewCurrentContext(config clientcmd.ClientConfig, newContext string) (clientcmd.ClientConfig, error) {
	rawConfig, err := config.RawConfig()
	if err != nil {
		panic(err)
	}

	detectedContext := false
	for name := range rawConfig.Contexts {
		if name == newContext {
			rawConfig.CurrentContext = newContext
			detectedContext = true
			break
		}
	}

	if !detectedContext {
		WrongContextMessage := fmt.Sprintf("context %v does not exist ", newContext)
		return config, errors.New(WrongContextMessage)
	}
	log.Debug("%v\n", config.ConfigAccess())
	err = clientcmd.ModifyConfig(config.ConfigAccess(), rawConfig, false)
	if err != nil {
		panic(err)
	}
	newConfig := clientcmd.NewDefaultClientConfig(rawConfig, &clientcmd.ConfigOverrides{})
	return newConfig, nil
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

func LoadConfig(opts ...string) (clientcmd.ClientConfig, error) {

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	rawConfig, err := loadingRules.Load()

	config := clientcmd.NewDefaultClientConfig(*rawConfig, &clientcmd.ConfigOverrides{})

	return config, err
}

//SetContextAction ....
func SetContextAction(c *cli.Context) error {

	var newContext string

	config, err := LoadConfig()
	if err != nil {
		return err
	}
	newContext = c.Args().First()
	oldContext := ""
	if newContext == "" {
		fmt.Println("missing context name ...")
		return errors.New("missing context name")
	}

	rawConfig, _ := config.RawConfig()
	oldContext = rawConfig.CurrentContext
	config, err = SetNewCurrentContext(config, newContext)

	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	if err != nil {
		fmt.Printf("wrong context %v \n", red(newContext))
		return err
	}

	if oldContext != newContext {
		fmt.Printf("switch context from %v to %v\n", oldContext, green(newContext))
	} else {
		fmt.Printf("context %v is already set as current\n", green(newContext))
	}
	return err
}
func HandleSetContext(c *cli.Context) error {

	var newContext string

	config, err := LoadConfig()
	if err != nil {
		return err
	}

	newContext = c.Args().Get(0)
	log.Debug("new context %v", newContext)

	if newContext != "" {
		fmt.Printf("updating context... %v", newContext)
		config, _ = SetNewCurrentContext(config, newContext)
	}
	rawConfig, err := config.RawConfig()
	fmt.Printf("\ncurrent context is %v\n", rawConfig.CurrentContext)

	flags := FlagOptions{Validate: false, Cache: c.GlobalBool("cache")}
	if flags.Cache {
		ListContextFromCache()
		return nil
	}
	ListContexts(config, flags)

	return nil
}
func HandleDeleteContext(c *cli.Context) error {

	red := color.New(color.FgRed).SprintFunc()

	config, err := LoadConfig()
	if err != nil {
		return err
	}

	rawConfig, err := config.RawConfig()
	fmt.Printf("\ncurrent context is %v\n", rawConfig.CurrentContext)

	if !c.Args().Present() {
		fmt.Println("no contexts to delete provided")
		return nil
	}

	args := c.Args()
	for i := 0; i < (len(c.Args().Tail()) + 1); i++ {
		arg := args.Get(i)
		_, ok := rawConfig.Contexts[arg]
		if !ok {
			fmt.Printf("nothing to delete - cannot find context %s\n", red(arg))
			continue
		}

		if rawConfig.CurrentContext == arg {
			fmt.Printf("warning: this removed your active context, use \"kubectl config use-context\" to select a different one\n")
			continue
		}

		delete(rawConfig.Contexts, arg)
		fmt.Printf("\ncontext %s deleted\n", red(arg))
	}

	if err := clientcmd.ModifyConfig(config.ConfigAccess(), rawConfig, true); err != nil {
		return err
	}

	return nil
}

//LoadConfig ...
func LoadConfigFromFile(opts ...string) (clientcmd.ClientConfig, error) {

	home, _ := os.UserHomeDir()
	var kubeconfig = filepath.Join(home, ".kube", "config")
	tempConfig, err := clientcmd.LoadFromFile(kubeconfig)
	if err != nil {
		return nil, err
	}
	green := color.New(color.FgGreen).SprintFunc()
	//tempConfig.CurrentContext
	config := clientcmd.NewDefaultClientConfig(*tempConfig, &clientcmd.ConfigOverrides{})
	fmt.Printf("%v\n", green(config.ConfigAccess().GetDefaultFilename()))

	return config, nil

}

//ListContextFromCache ...
func ListContextFromCache() {
	c, err := NewLocalCache()
	if err != nil {
		panic(err)
	}
	headers := []string{"Contexts", "IsAvailable", "AuthProvider"}

	printTable(headers, c.cache)
}

//UpdateCache
func UpdateCache(data map[string]*KubeContext) error {
	c, err := NewLocalCache()
	fmt.Printf("update cache %v\n", len(data))
	if err != nil {
		return err
	}
	for i, _ := range data {
		if data[i].Status == ClusterNotTested {
			_, ok := c.cache[i]
			if ok {
				data[i].Status = c.cache[i].Status
			}
		}
	}
	c.cache = data
	_, err = c.Flash()

	return err

}
func ListContexts(config clientcmd.ClientConfig, flags FlagOptions) {
	rawConfig, err := config.RawConfig()
	if err != nil {
		panic(err)
	}

	data := make(map[string]*KubeContext, 100)
	headers := []string{"Contexts", "IsAvailable", "AuthProvider"}

	for name, _ := range rawConfig.Contexts {
		var pods int
		pods = -1

		currentContextName := fmt.Sprintf("%v", name)

		cfg, err := clientcmd.NewDefaultClientConfig(rawConfig, &clientcmd.ConfigOverrides{CurrentContext: name}).ClientConfig()
		authProvider := "default"

		if err != nil {
			fmt.Printf("error is %v\n", err)
			authProvider = "invalid"
		}

		if (err == nil) && (cfg.AuthProvider != nil) {
			authProvider = cfg.AuthProvider.Name
		}
		log.Debug("[auth %v]\n", authProvider)
		//	auth := fmt.Sprintf("[%s]", authProvider)
		pods = 1
		var status ClusterStatus
		status = ClusterNotTested

		if flags.Validate {
			pods = testCluster(config, name)
			if pods != -1 {

				status = ClusterAvailable
			} else {
				status = ClusterNotAvailable
			}
		}
		namespace, _, _ := config.Namespace()
		isCurrentContext := (rawConfig.CurrentContext == name)
		kubeContext := &KubeContext{currentContextName, namespace,
			status, authProvider, time.Now().String(), isCurrentContext}
		data[currentContextName] = kubeContext

	}
	err = UpdateCache(data)
	if err != nil {
		panic(err)
	}
	printTable(headers, data)

}

//Contexts
type Contexts []string

func (a Contexts) Len() int { return len(a) }
func (a Contexts) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}
func (a Contexts) Less(i, j int) bool {

	return a[i] < a[j]
}

func printTable(header []string, data map[string]*KubeContext) {
	//happyIcon := "\u2714"
	//sadIcon := "\u2716"
	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	var names []string
	for name, _ := range data {
		names = append(names, name)
	}
	sort.Sort(Contexts(names))

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(header)

	for _, name := range names {
		ctx := data[name]
		var statusStr string
		contextName := ctx.Name
		switch ctx.Status {
		case ClusterAvailable:
			statusStr = fmt.Sprintf("%s", green("Yes"))
		case ClusterNotAvailable:
			statusStr = fmt.Sprintf("%s", red("No"))
		case ClusterNotTested:
			statusStr = "N/A"
		}
		if ctx.CurrentContext {
			contextName = green(contextName)
		}
		v := []string{contextName, statusStr, ctx.AuthProvider}
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

	log.Debug("[%v] pods are %v len=%v \n", currentContext, pods.Items, len(pods.Items))
	if err != nil {
		return -1
	}

	return len(pods.Items)

}

type FlagOptions struct {
	List     bool
	Validate bool
	Cache    bool
}

func CreateNamespace(ns string, config clientcmd.ClientConfig) error {
	c, err := config.ClientConfig()
	if err != nil {
		return err
	}
	restClient, err := kubernetes.NewForConfig(c)
	if err != nil {
		return err
	}
	nsSpec := &v1.Namespace{ObjectMeta: metav1.ObjectMeta{Name: ns}}

	_, err = restClient.CoreV1().Namespaces().Create(nsSpec)
	return err
}

//DeleteNamespace
func DeleteNamespace(ns string, config clientcmd.ClientConfig) error {
	c, err := config.ClientConfig()
	if err != nil {
		return err
	}
	restClient, err := kubernetes.NewForConfig(c)
	if err != nil {
		return err
	}

	err = restClient.CoreV1().Namespaces().Delete(ns, &metav1.DeleteOptions{})
	return err
}

//DeleteContexts
func DeleteContexts(contexts []string, config clientcmd.ClientConfig, modifyFile bool) error {

	rawConfig, err := config.RawConfig()
	if err != nil {

	}
	for _, ctx := range contexts {

		_, ok := rawConfig.Contexts[ctx]
		if !ok {
			fmt.Printf("nothing to delete - cannot find context %s\n", ctx)
			continue
		}

		delete(rawConfig.Contexts, ctx)
		fmt.Printf("\ncontext %s deleted\n", ctx)
	}
	if !modifyFile {
		return err
	}

	if err := clientcmd.ModifyConfig(config.ConfigAccess(), rawConfig, true); err != nil {
		return err
	}

	return nil

}

func SetNamespaceToContext(ns string, config clientcmd.ClientConfig) error {

	rawConfig, err := config.RawConfig()
	if err != nil {
		return err
	}
	currentCtxName := rawConfig.CurrentContext
	context := rawConfig.Contexts[currentCtxName]
	context.Namespace = ns

	err = clientcmd.ModifyConfig(config.ConfigAccess(), rawConfig, false)

	return err
}
