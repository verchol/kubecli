package config

import (
	"errors"
	"fmt"
	"strings"

	"log"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
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
		return ServiceAccount{}, err
	}

	restClient, err := kubernetes.NewForConfig(c)
	if err != nil {
		return ServiceAccount{}, err
	}
	spec := &v1.ServiceAccount{}
	spec.Name = name
	createdSa, err := restClient.CoreV1().ServiceAccounts(namespace).Create(spec)
	if err != nil {
		return ServiceAccount{}, err
	}
	watchOpts := metav1.ListOptions{}
	watcher, err := restClient.CoreV1().ServiceAccounts(namespace).Watch(watchOpts)
	if err != nil {
		return ServiceAccount{}, err
	}
	skip := false
	if len(createdSa.Secrets) != 0 {
		skip = true
	}
	eventChan := watcher.ResultChan()
	var sa *v1.ServiceAccount
	for !skip {

		select {
		case event := <-eventChan:
			log.Printf("event %v %v\n", event.Type, event.Object)
			sa = event.Object.(*v1.ServiceAccount)
			if (event.Type == watch.Modified || event.Type == watch.Added) && (len(sa.Secrets) > 0) {
				log.Printf("stop the loop \n")
				skip = true
				break
			}

		}

	}
	log.Printf("continue after loop\n")

	secrets := sa.Secrets
	if len(secrets) == 0 {
		return ServiceAccount{}, errors.New("no secret associated with service accoutn")
	}
	s := secrets[0]
	token := getSecretToken(restClient, namespace, s.Name)
	log.Printf("secret %v \n token %v\n", s.Name, token)

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
	log.Printf("sa = %v", sa.Name)
	secrets := sa.Secrets
	for _, s := range secrets {

		token := getSecretToken(restClient, namespace, s.Name)
		log.Printf("secret %v \n token %v\n", s.Name, token)

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
	log.Printf("\ncurrent context is %v\n", rawConfig.CurrentContext)

	context := currentCtx.DeepCopy()
	context.Namespace = namespace
	if satoken != "" {
		auth := clientcmdapi.NewAuthInfo()
		auth.Token = satoken
		authName := fmt.Sprintf("%v_%v_user", context.Cluster, namespace)
		log.Printf("user:%v", authName)
		rawConfig.AuthInfos[authName] = auth
		context.AuthInfo = authName
	}

	rawConfig.Contexts[contextName] = context
	rawConfig.CurrentContext = contextName

	err = clientcmd.ModifyConfig(config.ConfigAccess(), rawConfig, true)

	return err
}

//CreateNonAdminContext ...
func CreateNonAdminContext(contextName string, namespace string, config clientcmd.ClientConfig) error {

	sa := fmt.Sprintf("sa%v%v", namespace, contextName)
	sa = strings.ToLower(sa)
	roleName := fmt.Sprintf("role-%v-%v", contextName, namespace)
	roleName = strings.ToLower(roleName)

	roleOpts := NewRoleOpts(roleName, namespace)

	CreateNamespace(namespace, config)
	role, err := CreateRole(roleOpts, config)

	if err != nil {

		//TODO identify when the error reason is "AlreadyExists"
		//For now is skipping treating error

		//log.Fatal(err)
		//return err

	}

	saObj, err := CreateServiceAccount(namespace, sa, config)

	if err != nil {
		//TODO identify when the error reason is "AlreadyExists"
		//For now is skipping treating error

		log.Fatal(err)
		return err
	}

	roleBindingOpts := NewRoleBindingOpts(fmt.Sprintf("rb1-%v-%v", sa, namespace), namespace)
	roleBindingOpts.Role = role.Name
	roleBindingOpts.ServiceAccount = sa
	roleBindingOpts.ServiceAccountNs = namespace

	_, err = CreateRoleBinding(roleBindingOpts, config)
	if err != nil {
		//TODO identify when the error reason is "AlreadyExists"
		//For now is skipping treating error
		//log.Fatal(err)
	}

	err = CreateContext(contextName, namespace, string(saObj.Token), config)

	return err
}

//CreateAdminContext ...
func CreateAdminContext(contextToCreate string, namespace string, config clientcmd.ClientConfig) error {

	//consider add session number in the future

	serviceAccountName := fmt.Sprintf("sa%v%v", namespace, contextToCreate)
	serviceAccountName = strings.ToLower(serviceAccountName)
	clusterRoleName := fmt.Sprintf("clusterrole_%v_%v", namespace, contextToCreate)
	clusterRoleName = strings.ToLower(clusterRoleName)

	CreateNamespace(namespace, config)

	sa, err := CreateServiceAccount(namespace, serviceAccountName, config)
	if err != nil {
		//TODO identify already exists
		//return err
	}

	clusterRole, err := NewClusterRole(clusterRoleName, config)

	if err != nil {
		//TODO identify already exists
		//return err
	}

	roleBindingOpts := NewRoleBindingOpts(fmt.Sprintf("rb1-%v", serviceAccountName), namespace)
	roleBindingOpts.Role = clusterRole.Name
	roleBindingOpts.ServiceAccount = serviceAccountName
	roleBindingOpts.ServiceAccountNs = namespace

	_, err = CreateClusterRoleBinding(roleBindingOpts, config)
	if err != nil {
		//return err
	}
	err = CreateContext(contextToCreate, namespace, string(sa.Token), config)

	return err

}
