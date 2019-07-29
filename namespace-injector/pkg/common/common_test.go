package common

import (
	"fmt"
	"strings"
)

func ExampleListConfigMapCachedFiles() {
	WriteFile2Cache("test1_1", "aaa")
	WriteFile2Cache("test1_2", "aaa")
	WriteFile2Cache("test2_3", "aaa")
	WriteFile2Cache("test3_4", "aaa")
	res, _ := ListConfigMapCachedFiles("test1")
	fmt.Println(strings.Join(res, ","))
	RemoveFiles("test1_1", "test1_2", "test2_3", "test3_4")
	//Output: /tmp/test1_1,/tmp/test1_2
}
