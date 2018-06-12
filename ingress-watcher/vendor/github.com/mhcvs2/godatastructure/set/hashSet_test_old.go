package set

import (
	"fmt"
	"testing"
)

//样本测试
func ExampleHashSet_Add() {
	exampleSet := NewHashSet()
	fmt.Println(exampleSet.Add("a"))
	fmt.Println(exampleSet.String())
	fmt.Println(exampleSet.Add("a"))
	//Output:
	//true
	//Set{a}
	//false
}

func ExampleHashSet_Remove() {
	exampleSet := NewHashSet()
	exampleSet.Add("a")
	exampleSet.Add("b")
	exampleSet.Remove("b")
	exampleSet.Remove("b")
	fmt.Println(exampleSet.String())
	//Output: Set{a}
}

func ExampleHashSet_Clear() {
	exampleSet := NewHashSet()
	exampleSet.Add("a")
	exampleSet.Add("b")
	exampleSet.Clear()
	fmt.Println(exampleSet.String())
	//Output: Set{}
}

func ExampleHashSet_Contains() {
	exampleSet := NewHashSet()
	exampleSet.Add("a")
	fmt.Println(exampleSet.Contains("a"))
	//Output: true
}

func ExampleHashSet_Contains2() {
	exampleSet := NewHashSet()
	exampleSet.Add("a")
	fmt.Println(exampleSet.Contains("b"))
	//Output: false
}

func ExampleHashSet_Len() {
	exampleSet := NewHashSet()
	exampleSet.Add("a")
	fmt.Println(exampleSet.Len())
	exampleSet.Add("a")
	fmt.Println(exampleSet.Len())
	exampleSet.Add("b")
	fmt.Println(exampleSet.Len())
	//Output:
	// 1
	// 1
	// 2
}

func ExampleHashSet_Same() {
	exampleSet := NewHashSet()
	exampleSet.Add("a")
	exampleSet2 := NewHashSet()
	exampleSet2.Add("a")
	fmt.Println(exampleSet.Same(exampleSet2))
	exampleSet2.Add("b")
	fmt.Println(exampleSet.Same(exampleSet2))
	//Output:
	// true
	// false
}

func ExampleHashSet_Elements() {
	exampleSet := NewHashSet()
	exampleSet.Add("a")
	fmt.Println(exampleSet.Elements())
	//Output:
	//[a]
}

func ExampleHashSet_Copy() {
	exampleSet := NewHashSet()
	exampleSet.Add("a")
	var cp Set
	cp = exampleSet.Copy()
	fmt.Printf("set: %v\n", exampleSet.Elements())
	fmt.Printf("set_copy: %v\n", cp.Elements())
	exampleSet.Clear()
	fmt.Printf("set: %v\n", exampleSet.Elements())
	fmt.Printf("set_copy: %v\n", cp.Elements())
	//Output:
	//set: [a]
	//set_copy: [a]
	//set: []
	//set_copy: [a]
}

//功能测试
func TestHashSet_Elements(t *testing.T) {
	exampleSet := NewHashSet()
	exampleSet.Add("a")
	exampleSet.Add("b")
	exampleSet.Add("c")
	exampleSet.Add("d")
	for _, item := range exampleSet.Elements() {
		t.Log(item)
	}
}

//基准测试
func BenchmarkHashSet_Elements(b *testing.B) {
	exampleSet := NewHashSet()
	exampleSet.Add("a")
	exampleSet.Add("b")
	exampleSet.Add("c")
	exampleSet.Add("d")
	for _, item := range exampleSet.Elements() {
		b.Log(item)
	}
	b.StopTimer()
}
