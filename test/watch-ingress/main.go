package main

import (
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"fmt"
	"flag"
	"path/filepath"
	"os"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/informers"
	"time"
	extensions "k8s.io/api/extensions/v1beta1"
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}

func main() {
	ns := flag.String("ns", "kube-system", "namespace name")
	flag.Parse()
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

	resyncPeriod := 5 * time.Second

	infFactory := informers.NewFilteredSharedInformerFactory(clientset, resyncPeriod, *ns, func(*metav1.ListOptions) {})
	ingress := infFactory.Extensions().V1beta1().Ingresses().Informer()
	store := ingress.GetStore()
	is := store.List()
	fmt.Println(len(is))
	for _, i := range is {
		j := i.(*extensions.Ingress)
		fmt.Println(j.Name)
	}
}
