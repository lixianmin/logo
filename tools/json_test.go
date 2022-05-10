package tools

import (
	"fmt"
	"github.com/lixianmin/got/convert"
	"testing"
)

/********************************************************************
created:    2021-11-24
author:     lixianmin

*********************************************************************/

type stringStruct struct {
	age int
}

func (my *stringStruct) String() string {
	return fmt.Sprintf("age=%d", my.age)
}

func TestAppendJson(t *testing.T) {
	var s = &stringStruct{age: 37}
	var b []byte
	b = AppendJson(b, s)
	var result = convert.String(b)
	fmt.Println(result)
}
