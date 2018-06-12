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
	"k8s.io/client-go/tools/cache"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/api/core/v1"
	"os/signal"
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
	podInformer := sharedInformerFactory.Core().V1().Pods().Informer()
	go podInformer.Run(stopCh)

	if !cache.WaitForCacheSync(stopCh,
		podInformer.HasSynced,
	) {
		runtime.HandleError(fmt.Errorf("Timed out waiting for caches to sync"))
	} else {
		fmt.Println("success sync")
	}

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

	//pods2, err := podLister.Pods("kube-system").List(labels.Everything())
	//if err != nil {
	//	panic(err.Error())
	//}
	//fmt.Println(len(pods2))
	//for _, p := range pods2 {
	//	fmt.Printf("%s ", p.Name)
	//}
	//fmt.Println()
	//p, err := podLister.Pods("kube-system").Get("haproxy-k8s-ceph5")
	//if err != nil {
	//	panic(err.Error())
	//}
	//fmt.Println(p.Name)

	podEventHandler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			fmt.Println("create pod: ", pod.Name)
		},
		DeleteFunc: func(obj interface{}) {
			pod := obj.(*v1.Pod)
			fmt.Println("delete pod: ", pod.Name)
		},
		UpdateFunc: func(old, cur interface{}) {
			oldPod := old.(*v1.Pod)
			curPod := cur.(*v1.Pod)
			fmt.Println("update pod: ", oldPod.Name)
			fmt.Println(curPod.Name)
		},
	}
	podInformer.AddEventHandler(podEventHandler)

	select {
	case <-stopCh:
		os.Exit(0)
	case <-signalCh:
		fmt.Println("exit by signal")
		os.Exit(0)
	}
}
