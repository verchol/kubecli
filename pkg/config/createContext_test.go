package config

import (
	"fmt"
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const token = "ZXlKaGJHY2lPaUpTVXpJMU5pSXNJbXRwWkNJNklpSjkuZXlKcGMzTWlPaUpyZFdKbGNtNWxkR1Z6TDNObGNuWnBZMlZoWTJOdmRXNTBJaXdpYTNWaVpYSnVaWFJsY3k1cGJ5OXpaWEoyYVdObFlXTmpiM1Z1ZEM5dVlXMWxjM0JoWTJVaU9pSjBaWE4wTVNJc0ltdDFZbVZ5Ym1WMFpYTXVhVzh2YzJWeWRtbGpaV0ZqWTI5MWJuUXZjMlZqY21WMExtNWhiV1VpT2lKa1pXWmhkV3gwTFhSdmEyVnVMV2cxTlRod0lpd2lhM1ZpWlhKdVpYUmxjeTVwYnk5elpYSjJhV05sWVdOamIzVnVkQzl6WlhKMmFXTmxMV0ZqWTI5MWJuUXVibUZ0WlNJNkltUmxabUYxYkhRaUxDSnJkV0psY201bGRHVnpMbWx2TDNObGNuWnBZMlZoWTJOdmRXNTBMM05sY25acFkyVXRZV05qYjNWdWRDNTFhV1FpT2lKaE1HUTNNR0l4WlMwMlltTmpMVEV4WldFdFlUQXhaUzFqTmpNNFptRmhNamc0TTJNaUxDSnpkV0lpT2lKemVYTjBaVzA2YzJWeWRtbGpaV0ZqWTI5MWJuUTZkR1Z6ZERFNlpHVm1ZWFZzZENKOS52QUxoUHEySlNnRnFYZFYzaGJRM2ZCWVlSN0ZPa1RVWExzSkw2d3NNSG5QUGJwam80TVppZncxOUdsRy14SG1FaHVEZjlRSkY5REhlVk9fQ3M0eW5NNFFiY0FCMkplSmhQbXNYalVQenBIQjRORnJvRldXbWFaTm50UlNCQVFRdzQ2WVdsYzRzbFFrY21ZTkc2X3FVMmtCdjdjaUt2bzZNckRGTk9BUHg4MzFSV3N1N01MNGRfSHprNnFTdGtpZ1VJaVNUbm9Gb3BJdXhJeTI3OGJkYzVHWTEzaExlbkZHa0ZyVmUyczJISkVmcDZwaVFOUXJUVnVpYUF5LXdiQ0lrSDU3cnpxc3JxUHM1RHZQQkhnM3hRYUhMVXA3YjdGdHNtMTQ0VXl0d2pjN2ExRC0tbkNnM0NwQzJJQUxuLWIzWFRfVElFMVhNRUZROFBVdjFQZjFtYXlOVzZQVTlnaV9VZklHaGt4SXlydWFMRnUzd0l0azRrQVRmdFZMOU90NmIwMGFDZElqUkstVllmSnVweHlfX2kxTHAwWDhGYzhvYlFMdzFrWndmVlZzTVJRVXVIQ05uOVBDRkdSMWloVkZGMHRfSDNiQ3F4MEJEWjc4N1l0QlVIZ3dEVVpXbmtPUlNCYUdCcWdmNXUwdk1vWkhINFN1dm1qOTRVZGE4ZWg1VzFmb3F4ZlUzYXlwLS0xT09lSE03aFpHZVdzNVZtZnFXdnB4TU1Gb3EtNnBYM1ZYTllScVc0UTFHNWhUVFozaW5aeWxRQmQyVHRkeVRtSDR1R3lLb0VvMlF6NWRGbXRNNUZsV0dWRzlrRTZJY2NOeDhpYkR5aDZ6QV9hU0RhTnozY2RUVkpHZTNETV9TTVhsbUZzNTRiLWl2ai0zaGxTckFaRFRudWZJaTZBaw=="

func TestServiceAccount(t *testing.T) {
	//	sa, err := createServiceAccount("default", "sa3")
	sa, err := getServiceAccount("default", "default")
	fmt.Printf("%v %v", sa.Secrets, err)
}
func TestCreateServiceAccount(t *testing.T) {
	//	sa, err := createServiceAccount("default", "sa3")

	config, err := loadConfig()
	if err != nil {
		panic(err)
	}
	sa, err := createServiceAccount("test1", "testsa2", config)
	fmt.Printf("service account %v token- %v err- %v\n",
		sa.sa.Name, sa.token, err)
	//tokenBytes := make([]byte, len(sa.token))
	//_, err = base64.StdEncoding.Decode(tokenBytes, sa.token)
	//if err != nil {
	//	panic(err)
	//}
	role, err := createRole("test1", "testsa2-role", config)
	if err != nil {
		panic(err)
	}
	_, err = createRoleBinding("test1",
		"testsa2-rb", "test1",
		sa.sa.Name, role.Name, config)
	if err != nil {
		panic(err)
	}
	err = createContext("sa3ctx", "test1", string(sa.token), config)

	if err != nil {
		panic(err)
	}

}
func TestConnection(t *testing.T) {
	config, err := loadConfig()
	if err != nil {
		panic(err)
	}
	clientConfig, _ := config.ClientConfig()
	restClient, err := kubernetes.NewForConfig(clientConfig)
	if err != nil {
		panic(err)
	}

	pods, err := restClient.CoreV1().Pods("test1").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("[]There are %d pods in the cluster\n", len(pods.Items))
}
func TestCreateContext(t *testing.T) {
	config, err := loadConfig()
	ns := "test1"
	name := "ctx1"
	if err != nil {
		panic(err)
	}
	err = createContext(name, ns, token, config)
	if err != nil {
		panic(err)
	}

}

func TestCreateRole(t *testing.T) {
	config, err := loadConfig()
	ns := "test1"
	roleName := "myRole1"
	if err != nil {
		panic(err)
	}
	role, err := createRole(ns, roleName, config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("role %v with rules %v\n", role.ObjectMeta.Name, role.Rules)

}

func TestCreateRoleBinding(t *testing.T) {
	config, err := loadConfig()
	ns := "test1"
	sa := "sa1"
	saNamespace := "test1"
	roleName := "myRole1"
	roleBinding := "myrb1"
	if err != nil {
		panic(err)
	}
	rb, err := createRoleBinding(ns, roleBinding, saNamespace, sa, roleName, config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("roleBindings %v for role %vwith subjects %v\n", rb.ObjectMeta.Name, rb.RoleRef.Name, rb.Subjects)

}
