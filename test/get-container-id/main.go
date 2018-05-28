package main

import (
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"fmt"
	"flag"
	"path/filepath"
	myutils "k8s-dev/utils"
	"k8s.io/client-go/tools/clientcmd"
)


func main() {
	ns := flag.String("n", "kube-system", "namespace name")
	flag.Parse()

	kubeconfig := filepath.Join(myutils.HomeDir(), ".kube", "config")
	var err error
	var config *rest.Config
	if myutils.Exists(kubeconfig) {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	pods, err := clientset.CoreV1().Pods(*ns).List(metav1.ListOptions{LabelSelector:"k8s-app=kubernetes-dashboard"})
	for _, deploy := range pods.Items {
		fmt.Println(deploy.Name)
	}
}
