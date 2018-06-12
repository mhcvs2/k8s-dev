package main

import (
	"os/signal"
	"path/filepath"
	"os"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/kubernetes"
	log "github.com/sirupsen/logrus"
	"k8s-dev/ingress-watcher/watcher"
	"flag"
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}

func main() {
	debug := flag.Bool("debug", false, "debug")
	flag.Parse()
	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
	log.Debug("start debug...")
	signalCh := make(chan os.Signal, 1)
	go signal.Notify(signalCh, os.Interrupt, os.Kill)
	config, err := rest.InClusterConfig()
	if err != nil {
		home := homeDir()
		kubeconfig := filepath.Join(home, ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	stopCh := make(chan struct{})
	ingressWatcher := watcher.NewIngressWatcher(clientset)
	ingressWatcher.Run(stopCh)
	select {
	case <-stopCh:
		os.Exit(0)
	case <-signalCh:
		log.Info("exit by signal")
		os.Exit(0)
	}
}
