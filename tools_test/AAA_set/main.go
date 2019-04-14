package main

import (
	"fmt"
	"k8s.io/apimachinery/pkg/util/sets"
)

func t1() {
	stringSet := sets.NewString("a", "b")
	fmt.Println(stringSet.List())
	stringSet.Insert("haha")
	fmt.Println(stringSet.List())
	fmt.Println(stringSet.Has("haha"))
	fmt.Println(stringSet.HasAny("haha", "lala"))
	fmt.Println(stringSet.HasAll("haha", "a"))
}
//[a b]
//[a b haha]
//true
//true
//true

func t2() {
	testMap := make(map[string]int)
	testMap["aa"] = 2
	testMap["lala"] = 5
	keys := sets.StringKeySet(testMap)
	fmt.Println(keys.List())
	fmt.Println(keys.UnsortedList())
}
//[aa lala]

func t3() {

}

func main() {
	t2()
}
