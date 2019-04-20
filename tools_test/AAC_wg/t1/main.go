package main

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/wait"
	"time"
)

func Run1(stopCh <-chan struct{}) {
	fmt.Println("run-------------")
	<- stopCh
	fmt.Println("stop---------")
}

func t1() {
	stopCh := make(chan struct{})
	var wg wait.Group
	defer wg.Wait()
	wg.StartWithChannel(stopCh, Run1)
	time.Sleep(5 * time.Second)
	stopCh <- struct{}{}
}

func main() {
	t1()
}
