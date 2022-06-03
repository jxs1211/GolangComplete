package main

import (
	"fmt"
	"reflect"
	"runtime"
)

const s = `好雨知时节，当春乃发生。
随风潜入夜，润物细无声。
野径云俱黑，江船火独明。
晓看红湿处，花重锦官城。`

var i = int32(7)

func ShowCustomizeLiteralType() {
	i2 := int16(7)
	fmt.Printf("i: %d, type: %s\n", i, reflect.TypeOf(i).Kind().String())
	fmt.Printf("i2: %d, type: %s\n", i2, reflect.TypeOf(i2).Kind().String())
}

const (
	a        = 1
	b        = 2
	c int    = 1
	d myInt2 = 2
	e myInt  = 2
)

type myInt = int // type alias
type myInt2 int  // type redefine

func ShowUntypedConstant() {
	fmt.Println(a + b)
	fmt.Println(c + int(d))
	fmt.Println(c + e)
	fmt.Println(int(d) == e)
}

func main() {
	runtime.Breakpoint()
	// ShowCustomizeLiteralType()
	// ShowUntypedConstant()
	fmt.Println(s)
}
