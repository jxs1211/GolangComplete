package main

import (
	"fmt"
	"reflect"
)

type Status uint32

const (
	Open Status = iota
	Closed
)

func ShowNilInterface() {
	var p *int
	var i interface{} = p
	if i == nil {
		fmt.Println("i = nil")
		return
	}
	fmt.Printf("%v != nil, type: %T,\n i's value: %v, i's type: %s\n", i, i, reflect.ValueOf(i), reflect.TypeOf(i).String())
}

type MyStruct struct {
	Name     string
	nickname string
}

func ShowFuncListIteration() {
	type f func()

	f1 := func() {
		fmt.Println("call f1")
	}

	f2 := func() {
		fmt.Println("call f2")
	}
	funcs := []f{f1, f2}
	for _, v := range funcs {
		v()
	}

	funcs_m := map[string]f{
		"f1": f1,
		"f2": f2,
	}
	for k, v := range funcs_m {
		fmt.Println(k)
		v()
	}
	funcs_m["f2"]()
}

func main() {
	// fmt.Println("open: ", Open)
	// var v *net.IPAddr
	// fmt.Printf("v: %v\n", v)
	// var b bytes.Buffer
	// _, _ = b.Write([]byte("Hello"))
	// fmt.Println(b.String())

	// s := MyStruct{Name: "shen"}
	// fmt.Printf("%#v\n", s)

	// i := []int{'a': 1, 'b': 2}
	// fmt.Printf("i: %#v cap: %d char b: %d\n", i, cap(i), 'b')
	// m := map[string]int64{
	// 	"us": int64(time.Nanosecond),
	// 	// "us": int64(time.Nanosecond),
	// 	"µs": int64(time.Microsecond), // U+00B5 = micro symbol
	// 	"μs": int64(time.Microsecond),
	// }

	// fmt.Printf("%#v\n", m)
	// fmt.Printf("%q == %q: %v\n", "µs", "µs", "µs" == "µs")

	// ShowFuncListIteration()
}
