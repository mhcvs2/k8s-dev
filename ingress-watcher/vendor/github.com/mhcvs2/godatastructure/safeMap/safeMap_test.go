package safeMap

import (
	"testing"
	"fmt"
	"github.com/stretchr/testify/assert"
)

func TestSafeMap(t *testing.T) {
	sm := NewSafeMap()
	for i:=1; i<=10; i++{
		sm.Insert(fmt.Sprintf("key-%d", i), i)
	}
	a := assert.New(t)
	a.Equal(sm.Len(), 10)
	value, found := sm.Find("key-5")
	a.Equal(found, true)
	a.Equal(value, 5)
	sm.Delete("key-5")
	value, found = sm.Find("key-5")
	a.Equal(found, false)
	a.Equal(sm.Len(), 9)

	sm.Update("key-8", func(value interface{}, found bool) interface{} {
		if found{
			return 888
		}
		return 0
	})
	value, found = sm.Find("key-8")
	a.Equal(found, true)
	a.Equal(value, 888)

	m := sm.Close()
	fmt.Println(m)
}
