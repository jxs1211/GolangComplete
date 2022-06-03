package main

import "fmt"

type T struct {
	n int
	s string
}

func (T) M1() {}
func (T) M2() {}

type NonEmptyInterface interface {
	M1()
	M2()
}

func showInterface() {
	var t = T{
		n: 17,
		s: "hello, interface",
	}
	var ei interface{}
	ei = t

	var i NonEmptyInterface
	i = t
	fmt.Println(ei)
	fmt.Println(i)
}

type NonEmptyInterface2 interface {
	M1()
	M2()
}

type S struct {
	n string
	s int
}

func (S) M1() {}
func (S) M2() {}

func showInterface2() {
	s := S{n: "shen", s: 10}
	var ei interface{}
	var ii NonEmptyInterface2
	s.s = 20
	ei = s
	ii = s
	fmt.Println(ei)
	fmt.Println(ii)
}

func add(a, b int) int {
	return a + b
}
func main() {
	showInterface2()
}
