package set

import "testing"

func TestHashSet_Len_Contains(t *testing.T) {
	testSetLenAndContains(t, func() Set {return NewHashSet()}, "HashSet")
}
