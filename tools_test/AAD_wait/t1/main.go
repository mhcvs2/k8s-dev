package main

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/wait"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func t1() {
	run := func() {
		fmt.Println("run-----")
		fmt.Println(time.Now())
	}
	stopCh := make(chan struct{})
	wait.Until(run, 5 * time.Second, stopCh)
	signalCh := make(chan os.Signal)
	signal.Notify(signalCh, syscall.SIGTERM, os.Interrupt)
	<- signalCh
	stopCh <- struct{}{}
}
//run-----
//2019-04-22 11:08:09.500344057 +0800 CST m=+0.000493823
//run-----
//2019-04-22 11:08:14.500512847 +0800 CST m=+5.000662608
//run-----
//2019-04-22 11:08:19.50059721 +0800 CST m=+10.000746971
//run-----
//2019-04-22 11:08:24.500725132 +0800 CST m=+15.000874893

func main() {
	t1()
}
