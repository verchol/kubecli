package config

import (
	"fmt"
	"testing"

	"k8s.io/client-go/tools/clientcmd"
)

func TestCreateContext(t *testing.T) {

	config, err := loadConfig()
	if err != nil {
		panic(err)
	}

	rawConfig, err := config.RawConfig()
	currentCtx := rawConfig.Contexts[rawConfig.CurrentContext]
	fmt.Printf("\ncurrent context is %v\n", rawConfig.CurrentContext)

	context := currentCtx.DeepCopy()
	context.Namespace = "test1"

	rawConfig.Contexts["testcontext"] = context
	rawConfig.CurrentContext = "testcontext"

	if err := clientcmd.ModifyConfig(config.ConfigAccess(), rawConfig, true); err != nil {
		panic(err)
	}

}
