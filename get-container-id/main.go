package main

import (
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"fmt"
	"flag"
	"path/filepath"
	myutils "k8s-dev/utils"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"strings"
)


func main() {
	ns := flag.String("n", "default", "namespace name")
	label := flag.String("l", "", "namespace name")
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
	pods, err := clientset.CoreV1().Pods(*ns).List(lo)
	podNum := len(pods.Items)
	if podNum == 0 {
		fmt.Println("null")
		os.Exit(0)
	}
	if podNum > 1 {
		fmt.Println("multi")
		os.Exit(0)
	}
	pod := pods.Items[0]
	firstContainer := pod.Status.ContainerStatuses[0]
	fmt.Println(strings.TrimLeft(firstContainer.ContainerID,"docker://"))
}