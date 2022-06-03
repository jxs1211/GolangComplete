package main

import (
	"fmt"
	"reflect"
)

func showBoxing() {
	var n int = 61
	var ei interface{} = n
	n = 62
	fmt.Println("data in box:", ei)

	var m int = 51
	ei = &m
	m = 52

	p := ei.(*int)
	fmt.Println("data in box:", *p)
}

func showBoxing2() {
	var i interface{}
	a := 10 // boxing
	i = a   // value
	a = 20
	fmt.Println(i)

	b := 100
	i = &b
	b = 200 // pointer
	p := i.(*int)
	fmt.Println(*p)
}

func showType() {
	b := [...]byte{0x01, 0x02, 0x03}
	elemType := reflect.TypeOf(b)
	fmt.Println(elemType.Elem().Name())
	b = [3]byte{0x01, 0x02, 0x03}
	elemType = reflect.TypeOf(b)
	fmt.Println(elemType.Elem().Name())
	b2 := []byte{'1', '2', '3'}
	elemType = reflect.TypeOf(b2)
	fmt.Println(elemType.Elem().Name())
}

func main() {
	// showType()
	showBoxing2()
}
