package config

import (
	"errors"

	v1 "k8s.io/api/rbac/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

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
type RoleBindingOpts struct {
	Meta
	Role             string
	ServiceAccountNs string
	ServiceAccount   string
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
func NewRoleOpts(name string, ns string) *RoleOpts {

	r := RoleOpts{}
	r.Name = name
	r.Namespace = ns
	r.Verbs = []string{"create", "watch", "get", "list"}
	r.Resources = []string{"pods"}

	return &r

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

func CreateRole(opts *RoleOpts, config clientcmd.ClientConfig) (*v1.Role, error) {

	c, err := config.ClientConfig()
	if err != nil {
		panic(err)
	}
	restClient, err := kubernetes.NewForConfig(c)
	if err != nil {
		panic(err)
	}

	roleToCreate := v1.Role{}
	roleToCreate.ObjectMeta.Name = opts.Name
	roleToCreate.ObjectMeta.Namespace = opts.Namespace
	policy := v1.PolicyRule{Verbs: opts.Verbs, Resources: opts.Resources,
		APIGroups: []string{""}}
	roleToCreate.Rules = []v1.PolicyRule{policy}

	role, err := restClient.RbacV1().Roles(opts.Namespace).Create(&roleToCreate)

	if err != nil {
		panic(err.Error())
	}

	return role, err

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
		panic(err.Error())
	}

	return roleBinding, err

}
