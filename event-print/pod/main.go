package main

import (
	"flag"
	"fmt"
	myutils "k8s-dev/utils"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
)

func main() {
	ns := flag.String("n", "default", "namespace name")
	label := flag.String("l", "", "label name")
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
	var lo metav1.ListOptions
	if *label == "" {
		lo = metav1.ListOptions{}
	} else {
		lo = metav1.ListOptions{LabelSelector: *label}
	}
	watch, err := clientset.CoreV1().Pods(*ns).Watch(lo)
	if err != nil {
		panic(err)
	}

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	run := func() {
		fmt.Println("run...")
		for{
			select {
			case <- stopCh:
				fmt.Println("stop...")
				wg.Done()
			case event := <- watch.ResultChan():
				fmt.Printf("received event type: %v\n", event.Type)
				pod := event.Object.(*corev1.Pod)
				fmt.Printf("pod name is: %s\n", pod.Name)
			}
		}
	}

	go run()
}
