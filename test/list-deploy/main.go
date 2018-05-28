package main

import (
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"fmt"
)

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	deployments, err := clientset.AppsV1beta1().Deployments("default").List(metav1.ListOptions{})
	fmt.Println("deployments in namespace default:")
	for _, deploy := range deployments.Items {
		fmt.Println(deploy.Name)
	}
}
