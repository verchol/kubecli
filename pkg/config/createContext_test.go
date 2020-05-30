package config

import (
	"fmt"
	"testing"
	"time"

	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const token = "ZXlKaGJHY2lPaUpTVXpJMU5pSXNJbXRwWkNJNklpSjkuZXlKcGMzTWlPaUpyZFdKbGNtNWxkR1Z6TDNObGNuWnBZMlZoWTJOdmRXNTBJaXdpYTNWaVpYSnVaWFJsY3k1cGJ5OXpaWEoyYVdObFlXTmpiM1Z1ZEM5dVlXMWxjM0JoWTJVaU9pSjBaWE4wTVNJc0ltdDFZbVZ5Ym1WMFpYTXVhVzh2YzJWeWRtbGpaV0ZqWTI5MWJuUXZjMlZqY21WMExtNWhiV1VpT2lKa1pXWmhkV3gwTFhSdmEyVnVMV2cxTlRod0lpd2lhM1ZpWlhKdVpYUmxjeTVwYnk5elpYSjJhV05sWVdOamIzVnVkQzl6WlhKMmFXTmxMV0ZqWTI5MWJuUXVibUZ0WlNJNkltUmxabUYxYkhRaUxDSnJkV0psY201bGRHVnpMbWx2TDNObGNuWnBZMlZoWTJOdmRXNTBMM05sY25acFkyVXRZV05qYjNWdWRDNTFhV1FpT2lKaE1HUTNNR0l4WlMwMlltTmpMVEV4WldFdFlUQXhaUzFqTmpNNFptRmhNamc0TTJNaUxDSnpkV0lpT2lKemVYTjBaVzA2YzJWeWRtbGpaV0ZqWTI5MWJuUTZkR1Z6ZERFNlpHVm1ZWFZzZENKOS52QUxoUHEySlNnRnFYZFYzaGJRM2ZCWVlSN0ZPa1RVWExzSkw2d3NNSG5QUGJwam80TVppZncxOUdsRy14SG1FaHVEZjlRSkY5REhlVk9fQ3M0eW5NNFFiY0FCMkplSmhQbXNYalVQenBIQjRORnJvRldXbWFaTm50UlNCQVFRdzQ2WVdsYzRzbFFrY21ZTkc2X3FVMmtCdjdjaUt2bzZNckRGTk9BUHg4MzFSV3N1N01MNGRfSHprNnFTdGtpZ1VJaVNUbm9Gb3BJdXhJeTI3OGJkYzVHWTEzaExlbkZHa0ZyVmUyczJISkVmcDZwaVFOUXJUVnVpYUF5LXdiQ0lrSDU3cnpxc3JxUHM1RHZQQkhnM3hRYUhMVXA3YjdGdHNtMTQ0VXl0d2pjN2ExRC0tbkNnM0NwQzJJQUxuLWIzWFRfVElFMVhNRUZROFBVdjFQZjFtYXlOVzZQVTlnaV9VZklHaGt4SXlydWFMRnUzd0l0azRrQVRmdFZMOU90NmIwMGFDZElqUkstVllmSnVweHlfX2kxTHAwWDhGYzhvYlFMdzFrWndmVlZzTVJRVXVIQ05uOVBDRkdSMWloVkZGMHRfSDNiQ3F4MEJEWjc4N1l0QlVIZ3dEVVpXbmtPUlNCYUdCcWdmNXUwdk1vWkhINFN1dm1qOTRVZGE4ZWg1VzFmb3F4ZlUzYXlwLS0xT09lSE03aFpHZVdzNVZtZnFXdnB4TU1Gb3EtNnBYM1ZYTllScVc0UTFHNWhUVFozaW5aeWxRQmQyVHRkeVRtSDR1R3lLb0VvMlF6NWRGbXRNNUZsV0dWRzlrRTZJY2NOeDhpYkR5aDZ6QV9hU0RhTnozY2RUVkpHZTNETV9TTVhsbUZzNTRiLWl2ai0zaGxTckFaRFRudWZJaTZBaw=="

func TestLoadWithRules(t *testing.T) {

	//os.Setenv("KUBECONFIG", "/Users/verchol/dev/projects/kubecli/testdata/oketoconfig.yaml")
	config, err := LoadConfig()
	if err != nil {
		panic(err)
	}
	rawConfig, err := config.RawConfig()

	fmt.Printf("default context is %v \n", rawConfig.CurrentContext)
	contexts := rawConfig.Contexts

	for name, _ := range contexts {
		fmt.Printf("context is %v \n", name)
	}

}
func TestLocalCache(t *testing.T) {
	c, _ := NewLocalCache()

	config, err := LoadConfig()
	r, _ := config.RawConfig()

	context := r.Contexts[r.CurrentContext]
	kubeCtx := KubeContext{}
	kubeCtx.Name = r.CurrentContext
	kubeCtx.Namespace = context.Namespace
	kubeCtx.AuthProvider = context.AuthInfo

	_, err = c.AddEntry("t4", &kubeCtx).Flash()
	if err != nil {
		panic(err)
	}
	for k, v := range c.cache {
		fmt.Printf("cache  %v %v", k, v)
	}
	c.Reset()

}
func TestCluster(t *testing.T) {

	config, err := LoadConfig()
	if err != nil {
		fmt.Printf("can't load config")
		panic(err)
	}
	r, _ := config.RawConfig()
	contextToTest := r.CurrentContext
	fmt.Printf("context for validation %v \n", contextToTest)
	tempConfig := clientcmd.NewDefaultClientConfig(r,
		&clientcmd.ConfigOverrides{CurrentContext: contextToTest})

	namespace, _, err := tempConfig.Namespace()
	if err != nil {
		fmt.Printf("something wrong with context %v \n", contextToTest)
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
	before := time.Now()
	works, err := ValidateCluster(6, namespace, clientSet)
	d := time.Since(before)

	after := before.Add(d)

	fmt.Printf("started at %v\n", before)
	fmt.Printf("finieshed at %v\n", after)
	if !works || err != nil {
		t.Error(err)
	}
}
func TestRoleOpts(t *testing.T) {

	role := NewRoleOpts("role1", GlobalContext.Namespace)
	roleBinding := NewRoleBindingOpts("rb1", GlobalContext.Namespace)
	roleBinding.Role = role.Name
	roleBinding.ServiceAccount = "testsa1"
	roleBinding.ServiceAccountNs = GlobalContext.Namespace

}
func testServiceAccount(t *testing.T) {
	//	sa, err := createServiceAccount("default", "sa3")
	sa, err := getServiceAccount("default", "default")
	fmt.Printf("%v %v", sa.Secrets, err)
}
func TestCreateServiceAccount(t *testing.T) {
	//	sa, err := createServiceAccount("default", "sa3")

	config, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	namespace := "testsans"
	serviceAccountName := "testsa7"
	contextToTest := "sa3ctx"

	CreateNamespace(namespace, config)
	defer func() error {
		err := DeleteNamespace(namespace, config)
		return err
	}()
	sa, err := CreateServiceAccount(namespace, serviceAccountName, config)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Printf("service account %v token- %v err- %v\n",
		sa.Sa.Name, sa.Token, err)

	roleOpts := NewRoleOpts(fmt.Sprintf("role-%v", serviceAccountName), namespace)
	role, err := CreateRole(roleOpts, config)

	if err != nil {
		t.Fatal(err)
	}

	roleBindingOpts := NewRoleBindingOpts(fmt.Sprintf("rb1-%v", serviceAccountName), namespace)
	roleBindingOpts.Role = role.Name
	roleBindingOpts.ServiceAccount = serviceAccountName
	roleBindingOpts.ServiceAccountNs = namespace

	_, err = CreateRoleBinding(roleBindingOpts, config)
	if err != nil {
		t.Fatal(err)
	}
	err = CreateContext(contextToTest, namespace, string(sa.Token), config)

	if err != nil {
		t.Fatal(err)
	}

	//_, err = SetNewCurrentContext(config, restoreContext)

	if err != nil {
		t.Fatal(err)
	}

	err = DeleteContexts([]string{contextToTest}, config, true)

}
func TestConnection(t *testing.T) {
	config := GlobalContext.Config
	namespace := GlobalContext.Namespace

	clientConfig, _ := config.ClientConfig()
	restClient, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		panic(err)
	}

	pods, err := restClient.CoreV1().Pods(namespace).List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("[]There are %d pods in the cluster\n", len(pods.Items))
}
func TestCreateContext(t *testing.T) {
	config, err := LoadConfig()
	if err != nil {
		t.Fatal(err)
	}
	name := "ctx1"
	namespace := "testcreatecontexts"
	c, err := config.RawConfig()
	restoreContext := c.CurrentContext

	if err != nil {
		t.Fatal(err)
	}
	err = CreateContext(name, namespace, token, config)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		SetNewCurrentContext(config, restoreContext)
		DeleteNamespace(namespace, config)
		DeleteContexts([]string{name}, config, true)
	}()

}

func TestClusterRoleCreate(t *testing.T) {
	config := GlobalContext.Config
	_, err := CreateClusterRoleLogic(config)

	if err != nil {
		t.Error(err)
	}
	t.Log("succesfully created")
	DeleteClusterRoleLogic()
}
func TestCreateRoleWithoutDelete(t *testing.T) {

	_, err := CreateRoleLogic("role1", GlobalContext.Namespace, GlobalContext.Config)
	if err != nil {
		t.Error(err)
	}
}
func TestCreateRole(t *testing.T) {

	config := GlobalContext.Config
	opts, err := CreateRoleLogic("role1", GlobalContext.Namespace, GlobalContext.Config)
	if err != nil {
		t.Log(err)
		return
	}
	t.Logf("role created %v %v succesfully ", opts.Name, opts.Namespace)
	err = DeleteRole(opts, config)
	if err != nil {
		t.Log(err)
		return
	}

}
func TestCreateRoleBinding(t *testing.T) {
	config := GlobalContext.Config

	roleBindingOpts := NewRoleBindingOpts("myrb1", GlobalContext.Namespace)
	roleBindingOpts.Role = "myRole1"
	roleBindingOpts.ServiceAccount = "sa1"
	roleBindingOpts.ServiceAccountNs = "test1"

	rb, err := CreateRoleBinding(roleBindingOpts, config)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("roleBindings %v for role %vwith subjects %v\n", rb.ObjectMeta.Name, rb.RoleRef.Name, rb.Subjects)

}

func TestCreateAdminContext(t *testing.T) {
	var config clientcmd.ClientConfig
	if GlobalContext.Config == nil {
		var err error
		config, err = LoadConfig()
		if err != nil {
			t.Error(err)
			return
		}
	} else {
		config = GlobalContext.Config
	}

	c, err := config.RawConfig()
	if err != nil {
		t.Error(err)
		return
	}
	restoreContext := c.CurrentContext

	contextName := "ctxadmin"
	namespace := "admin1"

	err = CreateAdminContext(contextName, namespace, config)

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		SetNewCurrentContext(config, restoreContext)
		DeleteNamespace(namespace, config)
		DeleteContexts([]string{contextName}, config, true)
	}()

}
func CreateRoleLogic(name string, ns string, config clientcmd.ClientConfig) (*RoleOpts, error) {

	roleOpts := NewRoleOpts(name, ns)

	role, err := CreateRole(roleOpts, config)

	if err != nil {
		return roleOpts, err
	}
	fmt.Printf("role %v with rules %v\n", role.ObjectMeta.Name, role.Rules)

	return roleOpts, nil
}
func CreateClusterRoleLogic(config clientcmd.ClientConfig) (*v1.ClusterRole, error) {

	role, err := NewDefaultClusterRole(config)

	if err != nil {
		panic(err)
	}
	fmt.Printf("role %v with rules %v\n", role.ObjectMeta.Name, role.Rules)

	return role, err
}

func DeleteClusterRoleLogic() {

	config := GlobalContext.Config
	DeleteAdminClusterRole(config)

}
func CreateTestAdminContext(contextToCreate string, namespace string, config clientcmd.ClientConfig) error {

	return CreateAdminContext(contextToCreate, namespace, config)

}
func CreateAdminServiceAccount() {

}
func DeleteAdminServiceAccount() {

}

type GlobalTestContext struct {
	Config    clientcmd.ClientConfig
	Namespace string
}

var GlobalContext GlobalTestContext

func SetupTestContext(testContextName string, ns string) error {
	config, err := LoadConfig()
	if err != nil {
		fmt.Printf("warning %v\n", err)
	}
	err = CreateAdminContext(testContextName, ns, config)
	if err != nil {
		fmt.Printf("warning %v\n", err)
	}
	err = CreateNamespace(ns, config)
	if err != nil {
		fmt.Printf("warning %v\n", err)
	}
	err = SetNamespaceToContext(ns, config)
	if err != nil {
		fmt.Printf("warning %v\n", err)
	}

	GlobalContext = GlobalTestContext{config, ns}
	return err
}

func DeleteTestContext(testContextName string, ns string) error {
	config := GlobalContext.Config
	fmt.Printf("End of test execution , deleting context %v\n", testContextName)

	err := DeleteNamespace(ns, config)
	if err != nil {
		return err
	}
	err = DeleteContexts([]string{testContextName}, config, false)

	return err
}
