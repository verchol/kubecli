package config

import (
	"errors"

	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

var KubecliAdminRole = "KubecliAdminRole"

type Meta struct {
	Namespace string
	Name      string
}
type RoleOpts struct {
	Meta
	Verbs     []string
	Resources []string
	ApiGroups []string
}
type ClusterRoleOpts struct {
	RoleOpts
}
type RoleBindingOpts struct {
	Meta
	Role             string
	ServiceAccountNs string
	ServiceAccount   string
}
type ClusterRoleBindingOpts struct {
	RoleBindingOpts
}

func (r *RoleOpts) validate() (*RoleOpts, error) {
	if r.Name == "" {
		return r, errors.New("missing parameter : Name")
	}
	if r.Namespace == "" {
		return r, errors.New("missing parameter : Namespace")
	}

	if len(r.Verbs) == 0 {
		return r, errors.New("missing parameter : Verbs")
	}
	if len(r.Resources) == 0 {
		return r, errors.New("missing parameter : Resources")
	}

	return r, nil

}
func DefaultClusterRoleOpt(name string) *v1.ClusterRole {

	verbs := []string{"create", "watch", "get", "list"}
	resources := []string{"*"}
	apiGroups := []string{"", "extensions", "apps"}

	return NewClusterRoleOpt(name, verbs, resources, apiGroups)

}
func NewClusterRoleOpt(name string, verbs []string, resources []string,
	apiGroups []string) *v1.ClusterRole {

	roleToCreate := v1.ClusterRole{}
	roleToCreate.ObjectMeta.Name = name
	//roleToCreate.ObjectMeta.Namespace = opts.Namespace
	policy := v1.PolicyRule{Verbs: verbs, Resources: resources,
		APIGroups: apiGroups}
	roleToCreate.Rules = []v1.PolicyRule{policy}

	return &roleToCreate
}

func NewRoleOpts(name string, ns string) *RoleOpts {

	r := &RoleOpts{}
	r.Name = name
	r.Namespace = ns
	r.Verbs = []string{"create", "watch", "get", "list"}
	r.Resources = []string{"pods", "deployments"}
	r.ApiGroups = []string{"", "extensions", "apps"}

	return r

}
func NewRoleBindingOpts(name string, ns string) *RoleBindingOpts {

	r := RoleBindingOpts{}
	r.Name = name
	r.Namespace = ns

	return &r

}

func (r *RoleBindingOpts) validate() (*RoleBindingOpts, error) {
	if r.Name == "" {
		return r, errors.New("missing parameter : Name")
	}
	if r.Namespace == "" {
		return r, errors.New("missing parameter : Namespace")
	}

	if r.Role == "" {
		return r, errors.New("missing parameter : Role")
	}
	if r.ServiceAccount == "" {
		return r, errors.New("missing parameter : ServiceAccount")
	}

	if r.ServiceAccountNs == "" {
		return r, errors.New("missing parameter : ServiceAccountNs")
	}

	return r, nil

}

// T ...
type roleOptsGen func(r *RoleOpts) *RoleOpts
type roleBindingOptsGen func(r *RoleBindingOpts) *RoleBindingOpts

//DeleteRole ...
func DeleteRole(opts *RoleOpts, config clientcmd.ClientConfig) error {
	c, err := config.ClientConfig()
	if err != nil {
		return err
	}

	restClient, err := kubernetes.NewForConfig(c)
	if err != nil {
		return err
	}
	dopts := metav1.DeleteOptions{}
	err = restClient.RbacV1().Roles(opts.Namespace).Delete(opts.Name, &dopts)

	return err
}

//NewDefaultClusterRole
func NewDefaultClusterRole(config clientcmd.ClientConfig) (*v1.ClusterRole, error) {

	r := DefaultClusterRoleOpt(KubecliAdminRole)
	clusterRole, err := CreateAdminRole(r, config)

	return clusterRole, err

}

//NewClusterRole ...
func NewClusterRole(clusterRoleName string, config clientcmd.ClientConfig) (*v1.ClusterRole, error) {

	r := DefaultClusterRoleOpt(clusterRoleName)
	clusterRole, err := CreateAdminRole(r, config)

	return clusterRole, err

}

//CreateAdminRole

func DeleteAdminClusterRole(config clientcmd.ClientConfig) error {
	c, err := config.ClientConfig()
	if err != nil {
		panic(err)
	}
	restClient, err := kubernetes.NewForConfig(c)
	if err != nil {
		return err
	}

	_, err = restClient.RbacV1().ClusterRoles().Get(KubecliAdminRole, metav1.GetOptions{})

	if err != nil {
		return err
	}

	err = restClient.RbacV1().ClusterRoles().Delete(KubecliAdminRole, &metav1.DeleteOptions{})

	return err
}
func CreateAdminRole(roleToCreate *v1.ClusterRole, config clientcmd.ClientConfig) (*v1.ClusterRole, error) {
	c, err := config.ClientConfig()
	if err != nil {
		panic(err)
	}
	restClient, err := kubernetes.NewForConfig(c)
	if err != nil {
		return nil, err
	}

	role, err := restClient.RbacV1().ClusterRoles().Create(roleToCreate)

	if err != nil {
		return roleToCreate, err
	}

	return role, nil

}
func CreateRole(opts *RoleOpts, config clientcmd.ClientConfig) (*v1.Role, error) {

	c, err := config.ClientConfig()
	if err != nil {
		panic(err)
	}
	restClient, err := kubernetes.NewForConfig(c)
	if err != nil {
		return nil, err
	}

	roleToCreate := v1.Role{}
	roleToCreate.ObjectMeta.Name = opts.Name
	roleToCreate.ObjectMeta.Namespace = opts.Namespace
	policy := v1.PolicyRule{Verbs: opts.Verbs, Resources: opts.Resources,
		APIGroups: opts.ApiGroups}
	roleToCreate.Rules = []v1.PolicyRule{policy}

	role, err := restClient.RbacV1().Roles(opts.Namespace).Create(&roleToCreate)

	if err != nil {
		return &roleToCreate, err
	}

	return role, nil

}

func CreateRoleBinding(opts *RoleBindingOpts, config clientcmd.ClientConfig) (*v1.RoleBinding, error) {

	c, err := config.ClientConfig()
	if err != nil {
		panic(err)
	}
	restClient, err := kubernetes.NewForConfig(c)
	if err != nil {
		panic(err)
	}

	rb := &v1.RoleBinding{}
	rb.ObjectMeta.Name = opts.Name
	rb.ObjectMeta.Namespace = opts.Namespace
	rb.RoleRef = v1.RoleRef{Kind: "Role", Name: opts.Role}
	rb.Subjects = []v1.Subject{{Kind: "ServiceAccount", Name: opts.ServiceAccount, Namespace: opts.ServiceAccountNs}}

	roleBinding, err := restClient.RbacV1().RoleBindings(opts.Namespace).Create(rb)

	if err != nil {
		return roleBinding, err
	}

	return roleBinding, err

}
func CreateClusterRoleBinding(opts *RoleBindingOpts, config clientcmd.ClientConfig) (*v1.ClusterRoleBinding, error) {

	c, err := config.ClientConfig()
	if err != nil {
		panic(err)
	}
	restClient, err := kubernetes.NewForConfig(c)
	if err != nil {
		panic(err)
	}

	rb := &v1.ClusterRoleBinding{}
	rb.ObjectMeta.Name = opts.Name
	rb.ObjectMeta.Namespace = opts.Namespace
	rb.RoleRef = v1.RoleRef{Kind: "ClusterRole", Name: opts.Role}
	rb.Subjects = []v1.Subject{{Kind: "ServiceAccount", Name: opts.ServiceAccount, Namespace: opts.ServiceAccountNs}}

	roleBinding, err := restClient.RbacV1().ClusterRoleBindings().Create(rb)

	if err != nil {
		panic(err.Error())
	}

	return roleBinding, err

}
