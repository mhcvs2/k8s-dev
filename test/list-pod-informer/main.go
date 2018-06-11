package main

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"fmt"
	"path/filepath"
	"os"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/informers"
	"time"
	"k8s.io/apimachinery/pkg/labels"
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}

func main() {
	config, err := rest.InClusterConfig()
	if err != nil {
		home := homeDir()
		kubeconfig := filepath.Join(home, ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	sharedInformerFactory := informers.NewSharedInformerFactory(clientset, time.Minute*10)

	stopCh := make(chan struct{})
	sharedInformerFactory.Start(stopCh)
	time.Sleep(time.Second * 5)
	podLister := sharedInformerFactory.Core().V1().Pods().Lister()
	pods1, err := podLister.List(labels.NewSelector())
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(len(pods1))
	for _, p := range pods1 {
		fmt.Printf("%s ", p.Name)
	}
	fmt.Println()

	pods2, err := podLister.Pods("kube-system").List(labels.Everything())
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(len(pods2))
	for _, p := range pods2 {
		fmt.Printf("%s ", p.Name)
	}
	fmt.Println()
	p, err := podLister.Pods("kube-system").Get("swagger-ddd5d766c-xrwnr")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(p.Name)
}
