package main

import (
	"io/ioutil"
	"k8s.io/client-go/kubernetes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"fmt"
	"flag"
	"path"
	"path/filepath"
	"os"
	"k8s.io/client-go/tools/clientcmd"
)

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE")
}

func FileExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func WriteFile(name, content string) error {
	data := []byte(content)
	if err := ioutil.WriteFile(name, data, 0644); err != nil {
		return err
	}
	fmt.Println("write file " + name + " success")
	return nil
}


func main() {
	ns := flag.String("ns", "default", "namespace name")
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
	configmap, err := clientset.CoreV1().ConfigMaps(*ns).Get("mhctest", metav1.GetOptions{})

	for k, v := range configmap.Data {
		fmt.Println("key: " + k)
		fmt.Println(v)
		if err := WriteFile(path.Join("/tmp", k), v); err != nil {
			panic(err)
		}
	}
}
