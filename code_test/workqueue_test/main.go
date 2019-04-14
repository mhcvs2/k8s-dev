package main

import (
	"fmt"
	"k8s.io/client-go/util/workqueue"
	"time"
)

// 如果已经get出一个obj，没有forget，同样的obj无法再次入列直到forget之后
func t1() {
	queue := workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "Networks")
	queue.AddRateLimited("hahaha")
	go func() {
		for{
			fmt.Println("for-----")
			obj, _ := queue.Get()
			fmt.Println("got----")
			fmt.Println(obj)
			queue.Forget("hahaha")
			queue.Done(obj)
			time.Sleep(time.Second * 1)
		}
	}()
	time.Sleep(time.Second * 1)
	queue.AddRateLimited("hahaha")
	time.Sleep(time.Second * 1)
	queue.AddRateLimited("4325")
	time.Sleep(time.Second * 1)
	queue.AddRateLimited("hahaha")
	time.Sleep(time.Second * 10)
	time.Sleep(time.Second * 10)
}


func main() {
	t1()
}
