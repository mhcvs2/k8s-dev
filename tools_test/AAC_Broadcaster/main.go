package main

import (
	"fmt"
	"k8s.io/apimachinery/pkg/watch"
	"time"
)

func t1() {
	bc := watch.NewBroadcaster(10, watch.WaitIfChannelFull)
	w := bc.Watch()
	rc := w.ResultChan()
	go func() {
		for e := range rc {
			fmt.Println("get event: ")
			fmt.Println(e.Type)
		}
	}()
	bc.Action(watch.Added, nil)
	time.Sleep(time.Second * 3)
}
//get event:
//ADDED

func main() {
	t1()
}
