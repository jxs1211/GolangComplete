package main

import (
	"fmt"
	"runtime"
)

func show() {
	// original string
	var s string = "hello"
	fmt.Println("original string:", s)

	// reslice it and try to change the original string
	sl := []byte(s)
	sl[0] = 't'
	fmt.Println("slice:", string(sl))
	fmt.Println("after reslice, the original string is:", string(s))
}

func show2() {
	str := "hello"
	fmt.Println("raw string: ", str)

	s := []byte(str)
	s[0] = 't'
	fmt.Println("slice: ", string(s))
	fmt.Println("after reslice, the original string is:", string(str))
}

func main() {
	show2()
	runtime.Breakpoint()
}
