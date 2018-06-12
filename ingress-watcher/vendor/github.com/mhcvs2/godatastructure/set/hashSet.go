//Go HashSet implemented via map
package set

import (
	"bytes"
	"fmt"
)

//HashSet implement Set-----------------------------
type HashSet struct {
	m map[interface{}]bool
}

//Generate new HashSet
func NewHashSet() *HashSet {
	return &HashSet{m: make(map[interface{}]bool)}
}

//Add an element
//If contain, return false else true
func (set *HashSet) Add(e interface{}) bool {
	if !set.m[e] {
		set.m[e] = true
		return true
	}
	return false
}

//Remove an element
func (set *HashSet) Remove(e interface{}) {
	delete(set.m, e)
}

//Clear all elements
func (set *HashSet) Clear() {
	set.m = make(map[interface{}]bool)
}

//Determine whether the element is included in the HashSet
func (set *HashSet) Contains(e interface{}) bool {
	return set.m[e]
}

//Return number of elements in HashSet
func (set *HashSet) Len() int {
	return len(set.m)
}

//Determine whether the HashSet is same to other
func (set *HashSet) Same(other Set) bool {
	if other == nil {
		return false
	}
	if set.Len() != other.Len() {
		return false
	}
	for key := range set.m {
		if !other.Contains(key) {
			return false
		}
	}
	return true
}

//Return an iterable slice including all elements
func (set *HashSet) Elements() []interface{} {
	initialLen := len(set.m)
	snapshot := make([]interface{}, initialLen)
	actualLen := 0
	for key := range set.m {
		if actualLen < initialLen {
			snapshot[actualLen] = key
		} else {
			snapshot = append(snapshot, key)
		}
		actualLen++
	}
	if actualLen < initialLen {
		snapshot = snapshot[:actualLen]
	}
	return snapshot
}

//Return a fromat string including all elements
func (set *HashSet) String() string {
	var buf bytes.Buffer
	buf.WriteString("Set{")
	first := true
	for key := range set.m {
		if first {
			first = false
		} else {
			buf.WriteString(" ")
		}
		buf.WriteString(fmt.Sprintf("%v", key))
	}
	buf.WriteString("}")
	return buf.String()
}

//Return a copy HashSet
func (set *HashSet) Copy() Set {
	copySet := NewHashSet()
	for key := range set.m {
		copySet.Add(key)
	}
	return copySet
}
