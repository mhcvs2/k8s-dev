package main

import (
	"fmt"
	myutils "k8s-dev/utils"
	"k8s.io/client-go/kubernetes"
	apps "k8s.io/api/apps/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/cache"
	"path/filepath"
	kubeinformers "k8s.io/client-go/informers"
	"time"
	"k8s-dev/k8s-controller-custom-resource/pkg/signals"
)

func main() {
	stopCh := signals.SetupSignalHandler()
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
	var kubeCli kubernetes.Interface
	kubeCli, err = kubernetes.NewForConfig(config)
	var kubeInformerFactory kubeinformers.SharedInformerFactory
	kubeInformerFactory = kubeinformers.NewSharedInformerFactoryWithOptions(kubeCli, time.Second*10)
	setInformer := kubeInformerFactory.Apps().V1().StatefulSets()
	setInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: func(old, cur interface{}) {
			curSet := cur.(*apps.StatefulSet)
			oldSet := old.(*apps.StatefulSet)
			setName := curSet.GetName()
			if curSet.ResourceVersion == oldSet.ResourceVersion {
				// Periodic resync will send update events for all known statefulsets.
				// Two different versions of the same statefulset will always have different RVs.
				return
			}
			fmt.Println(setName)
		},
	})
	kubeInformerFactory.Start(stopCh)
	fmt.Println("start...")
	<-stopCh
	fmt.Println("end...")
}
