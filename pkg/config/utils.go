package config

import (
	"github.com/docker/machine/libmachine/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

//ValidateCluster ...
func ValidateCluster(waitingPeriod int64, namespace string, client *kubernetes.Clientset) (bool, error) {

	var t int64
	t = int64(waitingPeriod)
	if t == 0 {
		t = 6
	}
	pods, err :=
		client.
			CoreV1().
			Pods(namespace).
			List(metav1.ListOptions{TimeoutSeconds: &t})

	log.Debug("[%v] pods are %v len=%v \n", pods.Items, len(pods.Items))
	if err != nil {
		return false, err
	}

	return true, nil

}
