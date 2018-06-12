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
	"k8s.io/client-go/tools/cache"
	"k8s.io/apimachinery/pkg/util/runtime"
	"os/signal"
	"k8s.io/api/extensions/v1beta1"
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}

func main() {
	signalCh := make(chan os.Signal, 1)
	go signal.Notify(signalCh, os.Interrupt, os.Kill)
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
	ingInformer := sharedInformerFactory.Extensions().V1beta1().Ingresses().Informer()
	go ingInformer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh,
		ingInformer.HasSynced,
	) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
	} else {
		fmt.Println("success sync")
	}

	ingEventHandler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			ing := obj.(*v1beta1.Ingress)
			fmt.Println("create ingress: ", ing.Name)
			rules := ing.Spec.Rules
			for _, rule := range rules {
				fmt.Println(rule.Host)
			}
			fmt.Println("------------------------")
		},
		DeleteFunc: func(obj interface{}) {
			ing := obj.(*v1beta1.Ingress)
			fmt.Println("delete ingress: ", ing.Name)
		},
		UpdateFunc: func(old, cur interface{}) {
			oldIng := old.(*v1beta1.Ingress)
			curIng := cur.(*v1beta1.Ingress)
			fmt.Println("update ingress: ", oldIng.Name)
			fmt.Println(curIng.Name)
		},
	}
	ingInformer.AddEventHandler(ingEventHandler)

	select {
	case <-stopCh:
		os.Exit(0)
	case <-signalCh:
		fmt.Println("exit by signal")
		os.Exit(0)
	}
}
