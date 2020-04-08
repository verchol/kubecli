package config

import (
	"errors"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type ServiceAccount struct {
	Sa    *v1.ServiceAccount
	Token []byte
}

//CreateServiceAccount ...
func CreateServiceAccount(namespace string, name string, config clientcmd.ClientConfig) (ServiceAccount, error) {

	c, err := config.ClientConfig()
	if err != nil {
		panic(err)
	}
	restClient, err := kubernetes.NewForConfig(c)
	if err != nil {
		panic(err)
	}
	spec := &v1.ServiceAccount{}
	spec.Name = name
	createdSa, err := restClient.CoreV1().ServiceAccounts(namespace).Create(spec)
	if err != nil {
		panic(err.Error())
	}
	sa, err := restClient.CoreV1().ServiceAccounts(namespace).Get(createdSa.Name, metav1.GetOptions{})
	fmt.Printf("sa = %v", sa)
	secrets := sa.Secrets
	if len(secrets) == 0 {
		panic(errors.New("no secrets associated with service account"))
	}
	s := secrets[0]

	if len(secrets) == 0 {
		panic(errors.New("no secretes associated with sa"))
	}

	token := getSecretToken(restClient, namespace, s.Name)
	fmt.Printf("secret %v \n token %v\n", s.Name, token)

	return ServiceAccount{sa, token}, err

}
func getSecretToken(client *kubernetes.Clientset, ns string, name string) []byte {
	secret, err := client.CoreV1().Secrets(ns).Get(name, metav1.GetOptions{})
	if err != nil {
		panic(err)
	}
	if secret.Type == "kubernetes.io/service-account-token" {
		token := secret.Data["token"]
		return token

	}

	return []byte{}

}
func getServiceAccount(namespace string, name string) (*v1.ServiceAccount, error) {
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	clientConfig, _ := config.ClientConfig()
	restClient, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		panic(err)
	}

	sa, err := restClient.CoreV1().ServiceAccounts(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("sa = %v", sa.Name)
	secrets := sa.Secrets
	for _, s := range secrets {

		token := getSecretToken(restClient, namespace, s.Name)
		fmt.Printf("secret %v \n token %v\n", s.Name, token)

	}
	return sa, err

}

//CreateContext ...
func CreateContext(contextName string, namespace string, satoken string, config clientcmd.ClientConfig) error {
	//choose cluster + namespace + user
	//create
	//set as default

	rawConfig, err := config.RawConfig()
	currentCtx := rawConfig.Contexts[rawConfig.CurrentContext]
	fmt.Printf("\ncurrent context is %v\n", rawConfig.CurrentContext)

	context := currentCtx.DeepCopy()
	context.Namespace = namespace
	if satoken != "" {
		auth := clientcmdapi.NewAuthInfo()
		auth.Token = satoken
		authName := fmt.Sprintf("%v_%v_user", context.Cluster, namespace)
		fmt.Printf("user:%v", authName)
		rawConfig.AuthInfos[authName] = auth
		context.AuthInfo = authName
	}

	rawConfig.Contexts[contextName] = context
	rawConfig.CurrentContext = contextName

	err = clientcmd.ModifyConfig(config.ConfigAccess(), rawConfig, true)

	return err
}
