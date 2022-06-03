package main

import "fmt"

type T struct {
	a int
}

func (t T) M1() {
	t.a = 10
}

func (t *T) M2() {
	t.a = 11
}

func showN1() {
	var t T // t.a = 0
	println(t.a)

	t.M1()
	println(t.a)

	t.M2()
	println(t.a)
}

type T2 struct {
	a int
}

func (t T2) Get() int {
	return t.a
}

func (t T2) Set(i int) int {
	t.a = i
	return t.a
}

func showN2() {
	var a T2
	fmt.Println(a.Get())
	fmt.Println(T2.Get(a))
	fmt.Println(a.Set(11))
	fmt.Println(T2.Set(a, 11))
}

type MyInt int

func (i MyInt) String() string {
	return fmt.Sprintf("%d", int(i))
}

// func (i int)String()string{
// 	return fmt.Sprintf("%d", int(i))
// }

func main() {
	showN2()
}
