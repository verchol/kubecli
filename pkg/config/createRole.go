package config

import (
	v1 "k8s.io/api/rbac/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func createRole(namespace string, roleName string, config clientcmd.ClientConfig) (*v1.Role, error) {

	c, err := config.ClientConfig()
	if err != nil {
		panic(err)
	}
	restClient, err := kubernetes.NewForConfig(c)
	if err != nil {
		panic(err)
	}

	roleToCreate := v1.Role{}
	roleToCreate.ObjectMeta.Name = roleName
	roleToCreate.ObjectMeta.Namespace = namespace
	policy := v1.PolicyRule{Verbs: []string{"get", "watch", "list"}, Resources: []string{"pods"}, APIGroups: []string{""}}
	roleToCreate.Rules = []v1.PolicyRule{policy}

	role, err := restClient.RbacV1().Roles(namespace).Create(&roleToCreate)

	if err != nil {
		panic(err.Error())
	}

	return role, err

}

func createRoleBinding(namespace string, name string,
	saNamespace string, serviceAccount string,
	role string, config clientcmd.ClientConfig) (*v1.RoleBinding, error) {

	c, err := config.ClientConfig()
	if err != nil {
		panic(err)
	}
	restClient, err := kubernetes.NewForConfig(c)
	if err != nil {
		panic(err)
	}

	rb := &v1.RoleBinding{}
	rb.ObjectMeta.Name = name
	rb.ObjectMeta.Namespace = namespace
	rb.RoleRef = v1.RoleRef{Kind: "Role", Name: role}
	rb.Subjects = []v1.Subject{{Kind: "ServiceAccount", Name: serviceAccount, Namespace: saNamespace}}

	roleBinding, err := restClient.RbacV1().RoleBindings(namespace).Create(rb)

	if err != nil {
		panic(err.Error())
	}

	return roleBinding, err

}
