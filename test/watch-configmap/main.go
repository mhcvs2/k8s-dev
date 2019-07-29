package main

import (
	"fmt"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"os/signal"
	"path/filepath"
	"time"
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

	//sharedInformerFactory := informers.NewSharedInformerFactory(clientset, time.Minute*10)
	sharedInformerFactory := informers.NewSharedInformerFactoryWithOptions(clientset, time.Minute*10, informers.WithNamespace("default"))
	stopCh := make(chan struct{})
	ingInformer := sharedInformerFactory.Core().V1().ConfigMaps().Informer()
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
			ing := obj.(*v1.ConfigMap)
			fmt.Println("create configmap: ", ing.Name)
			fmt.Println(ing.ObjectMeta.Annotations)
			fmt.Println("------------------------")
		},
		DeleteFunc: func(obj interface{}) {
			ing := obj.(*v1.ConfigMap)
			fmt.Println("delete configmap: ", ing.Name)
		},
		UpdateFunc: func(old, cur interface{}) {
			oldIng := old.(*v1.ConfigMap)
			curIng := cur.(*v1.ConfigMap)
			fmt.Println("update configmap: ", oldIng.Name)
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
