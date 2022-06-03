package main

import (
	"fmt"
	"io"
	"net/http"
)

type Worker interface {
	work()
}

type Person struct {
	name string
}

type Person2 struct {
	name string
	age  int
}

type Person3 struct {
	name string
	age  int
	Worker
}

func (p Person) work() {
	fmt.Println("Person work!")
}

func (p Person2) work() {
	fmt.Println("Person2 work!")
}

// func (p Person3) work() {
// 	fmt.Println("Person3 work!")
// }

func ShowCompareInterfaceType() {
	var w1 Worker = Person2{name: "shen"}
	var w2 Worker = Person3{name: "shen", Worker: w1}
	w2.work()
	// var w2 Worker = Person2{name: "shen", age: 10}

	// var w2 Worker = Person2{name: "shen"}
	if w1 != w2 {
		fmt.Printf("%v != %v\n", w1, w2)
		return
	}
	fmt.Printf("%v = %v\n", w1, w2)
}

func ShowCompareStructType() {
	p1 := Person2{name: "shen"}
	p2 := Person2{name: "shen", age: 10}
	// p2 := Person2{name: "shen"}
	if p1 != p2 {
		fmt.Printf("%v != %v\n", p1, p2)
		return
	}
	fmt.Printf("%v = %v\n", p1, p2)
}

func main() {
	// ShowCompareStructType()
	ShowCompareInterfaceType()
	var errr error
	var r io.Reader
	h := http.Handler
}
