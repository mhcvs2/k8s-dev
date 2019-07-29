package main

import (
	"flag"
	"github.com/Sirupsen/logrus"
	"k8s-dev/namespace-injector/pkg"
	"k8s-dev/namespace-injector/pkg/common"
	"k8s-dev/namespace-injector/pkg/config"
	"k8s-dev/namespace-injector/pkg/signals"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"path/filepath"
	"time"
)

func main() {
	flag.Parse()
	if err := common.MakeCacheDir(); err != nil {
		panic(err)
	}
	stopCh := signals.SetupSignalHandler()
	cfg, err := rest.InClusterConfig()
	if err != nil {
		home := common.HomeDir()
		kubeconfig := filepath.Join(home, ".kube", "config")
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}
	common.CheckKubeCtl()
	clientset, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		panic(err.Error())
	}
	informerFactory := informers.NewSharedInformerFactoryWithOptions(clientset, time.Second*30, informers.WithNamespace(config.ConfigMapNamespace))

	controller := pkg.NewController(
		clientset,
		informerFactory.Core().V1().ConfigMaps(),
		informerFactory.Core().V1().Namespaces(),
		)
	go informerFactory.Start(stopCh)
	if err = controller.Run(1, stopCh); err != nil {
		logrus.Fatalf("Error running controller: %s", err.Error())
	}

}