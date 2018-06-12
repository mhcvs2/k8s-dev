package set

import (
	"testing"
	"fmt"
)

func New() (s1, s2 *HashSet) {
	s1 = NewHashSet()
	s1.Add("a")
	s1.Add("b")
	s1.Add("c")

	s2 = NewHashSet()
	s2.Add("d")
	s2.Add("b")
	s2.Add("c")
	return s1, s2
}

func check(strs []string, set Set, t *testing.T) {
	for _, i := range strs {
		if !set.Contains(i) {
			t.Errorf("set union shoule contain %s", i)
		}
	}
}

func ExampleIsSuperset() {
	s1 := NewHashSet()
	s1.Add("a")
	fmt.Println(IsSuperset(s1, nil))
	fmt.Println(IsSuperset(nil, nil))
	s2 := NewHashSet()
	fmt.Println(IsSuperset(s1, s2))
	//Output:
	//false
	//false
	//true
}

func ExampleIsSuperset2() {
	s1, s2 := New()
	fmt.Println(IsSuperset(s1, s2))
	s1.Remove("a")
	s2.Remove("d")
	fmt.Println(IsSuperset(s1, s2))
	s2.Remove("c")
	fmt.Println(IsSuperset(s1, s2))
	//Output:
	//false
	//false
	//true
}

func TestUnion(t *testing.T) {
	s1, s2 := New()
	union := Union(s1, s2)
	t.Log(union.String())
	strs := []string{"a", "b", "c", "d"}
	check(strs, union, t)
}

func ExampleUnion() {
	s1 := NewHashSet()
	s1.Add("a")
	fmt.Println(Union(s1, nil))
	fmt.Println(Union(nil, nil))
	s2 := NewHashSet()
	fmt.Println(Union(s1, s2))
	//Output:
	// Set{a}
	//<nil>
	//Set{a}
}

func TestIntersect(t *testing.T) {
	s1, s2 := New()
	intersect := Intersect(s1, s2)
	t.Log(intersect.String())
	strs := []string{"b", "c"}
	check(strs, intersect, t)
}

func ExampleIntersect() {
	s1 := NewHashSet()
	s1.Add("a")
	fmt.Println(Intersect(s1, nil))
	fmt.Println(Intersect(nil, nil))
	s2 := NewHashSet()
	fmt.Println(Intersect(s1, s2))
	//Output:
	//<nil>
	//<nil>
	//<nil>
}

func ExampleDifference() {
	s1 := NewHashSet()
	s1.Add("a")
	fmt.Println(Difference(s1, nil))
	fmt.Println(Difference(nil, nil))
	s2 := NewHashSet()
	fmt.Println(Difference(s1, s2))
	//Output:
	//Set{a}
	//<nil>
	//Set{a}
}

func ExampleDifference2() {
	s1, s2 := New()
	fmt.Println(Difference(s1, s2))
	//Output:
	//Set{a}
}

func ExampleSummetricDifference() {
	s1 := NewHashSet()
	s1.Add("a")
	fmt.Println(Difference(s1, nil))
	fmt.Println(Difference(nil, nil))
	s2 := NewHashSet()
	fmt.Println(Difference(s1, s2))
	//Output:
	//Set{a}
	//<nil>
	//Set{a}
}

func TestSummetricDifference(t *testing.T) {
	s1, s2 := New()
	intersect := SummetricDifference(s1, s2)
	t.Log(intersect.String())
	strs := []string{"a", "d"}
	check(strs, intersect, t)
}