package main

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/tools/cache"
)

func t1() {
	indices := cache.Indices{}
	s := cache.NewThreadSafeStore(cache.Indexers{}, indices)
	indesers := cache.Indexers{}
	// 索引函数输入value，返回索引名
	// indices 存的是索引名到(key的集合)的map
	indesers["list"] = func(obj interface{}) (strings []string, e error) {
		return []string{"some"}, nil
	}
	s.AddIndexers(indesers)
	s.Add("a", "bb")
	s.Add("b", "cc")
	indice := cache.Index{}
	ss := sets.NewString("a")
	indice["some"] = ss
	indices["list"] = indice
	if v, ok := s.Get("a"); ok {
		fmt.Println(v)
	} else {
		fmt.Println("not exist")
	}

	if data, err := s.Index("list", "bbc"); err != nil {
		fmt.Println("err: " + err.Error())
	} else {
		fmt.Println(data)
	}
	if data, err := s.ByIndex("list", "some"); err != nil {
		fmt.Println("err: " + err.Error())
	} else {
		fmt.Println(data)
	}

	if data, err := s.IndexKeys("list", "some"); err != nil {
		fmt.Println("err: " + err.Error())
	} else {
		fmt.Println(data)
	}
}

func main() {
	t1()
}