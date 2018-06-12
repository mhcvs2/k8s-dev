package stack

import (
	"fmt"
	"container/list"
)

func printList(l *list.List) {
	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value)
	}
}

func ExampleStack_Push() {
	s := NewStack()
	s.Push(2)
	s.Push(5)
	s.Push(8)
	printList(s.list)
	//Output:
	//2
	//5
	//8
}

func ExampleStack_Pop() {
	s := NewStack()
	s.Push(2)
	s.Push(5)
	s.Push(8)
	fmt.Println(s.Pop())
	printList(s.list)
	//Output:
	//8
	//2
	//5
}

func ExampleStack_Peak() {
	s := NewStack()
	s.Push(2)
	s.Push(5)
	s.Push(8)
	fmt.Println(s.Peak())
	printList(s.list)
	//Output:
	//8
	//2
	//5
	//8
}

func ExampleStack_Len() {
	s := NewStack()
	s.Push(2)
	s.Push(5)
	s.Push(8)
	fmt.Println(s.Len())
	//Output:
	//3
}

func ExampleStack_Empty() {
	s := NewStack()
	fmt.Println(s.Empty())
	s.Push(2)
	s.Push(5)
	s.Push(8)
	fmt.Println(s.Empty())
	fmt.Println(s.Pop())
	fmt.Println(s.Pop())
	fmt.Println(s.Pop())
	fmt.Println(s.Empty())
	//Output:
	//true
	//false
	//8
	//5
	//2
	//true
}
