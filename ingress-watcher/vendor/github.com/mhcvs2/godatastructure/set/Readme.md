Example
==========================
<pre><code>
package main

import (
	"github.com/mhcvs2/godatastructure/set"
	"fmt"
)

func t1() {
	exampleSet := set.NewHashSet()
	exampleSet.Add("a")
	exampleSet.Add("b")
	exampleSet.Add("c")
	fmt.Println(exampleSet.String())
	exampleSet.Remove("a")
	exampleSet.Add("b")
	fmt.Println(exampleSet.String())
}
//Set{a b c}
//Set{b c}

func t2() {
	s1 := set.NewHashSet()
	s1.Add("a")
	s1.Add("b")
	s1.Add("c")

	s2 := set.NewHashSet()
	s2.Add("d")
	s2.Add("b")
	s2.Add("c")

	fmt.Println(set.Union(s1, s2).String())
	fmt.Println(set.Intersect(s1, s2).String())
	fmt.Println(set.Difference(s1, s2).String())
	fmt.Println(set.SummetricDifference(s1, s2).String())
}

//Set{c d b a}
//Set{b c}
//Set{a}
//Set{d a}
</pre></code>