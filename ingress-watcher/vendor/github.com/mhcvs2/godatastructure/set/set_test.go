package set

import (
	"testing"
	"math/rand"
	"time"
	"bytes"
)

func testSetLenAndContains(t *testing.T, newSet func() Set, typeName string) {
	t.Logf("Starting Test%sLenAndContains...", typeName)
	set, expectedElemMap := genRandSet(newSet)
	t.Logf("Got a %s value: %v.", typeName, set)
	expectedLen := len(expectedElemMap)
	if set.Len() != expectedLen {
		t.Errorf("ERROR： The length of %s value %d is not %d!\n",
			set.Len(), typeName, expectedLen)
		t.FailNow()
	}
	t.Logf("The length of %s value is %d.\n", typeName, set.Len())
	for k := range expectedElemMap {
		if !set.Contains(k) {
			t.Errorf("ERROR: The %s value %v do not contains %v!",
				set, typeName, k)
			t.FailNow()
		}
	}
}


//--gen rand function--
func genRandSet(newSet func() Set) (set Set, elemMap map[interface{}]bool) {
	set = newSet()
	elemMap = make(map[interface{}]bool)
	var enouth bool
	for !enouth {
		e := genRandElement()
		set.Add(e)
		elemMap[e] = true
		if len(elemMap) >= 3 {
			enouth = true
		}
	}
	return
}

func genRandElement() interface{} {
	seed := rand.Int63n(10000)
	switch seed {
	case 0:
		return genRandInt()
	case 1:
		return genRandString()
	case 2:
		return struct {
			num int64
			str string
		}{genRandInt(), genRandString()}
	default:
		const length = 2
		arr := new([length]interface{})
		for i:=0; i<length; i++ {
			if i%2 ==0 {
				arr[i] = genRandInt()
			} else {
				arr[i] = genRandString()
			}
		}
		return *arr
	}
}

func genRandString() string {
	var buff bytes.Buffer
	var prev string
	var curr string
	for i := 0; buff.Len() < 3; i++ {
		curr = string(genRandAZAscii())
		if curr == prev {
			continue
		} else {
			prev = curr
		}
		buff.WriteString(curr)
	}
	return buff.String()
}

func genRandAZAscii() int {
	min := 65
	max := 90
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}

func genRandInt() int64 {
	return rand.Int63n(10000)
}