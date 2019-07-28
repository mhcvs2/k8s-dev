package main

import (
	"fmt"
	"k8s-dev/namespace-injector/config"
)

func main() {
	fmt.Println(config.NamespaceLabelSelectors)
}