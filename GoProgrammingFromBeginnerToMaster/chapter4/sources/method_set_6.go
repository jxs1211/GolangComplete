package main

import "fmt"

type Interface interface {
	M1()
	M2()
}

type T struct {
	Interface
}

func (T) M1() {
	println("T's M1")
}

type S struct{}

func (S) M1() {
	println("S's M1")
}
func (S) M2() {
	println("S's M2")
}

func showStructEmbedInterface() {
	var t = T{
		Interface: S{},
	}

	t.M1()
	t.M2()
}

type Interface2 interface {
	M1()
	M2()
}

type Interface3 interface {
	M2()
}

type T2 struct {
	Interface2
	// Interface3
}

func (T2) M1() {
	fmt.Println("T2 M1")
}

type S2 struct {
}

func (S2) M1() {
	fmt.Println("S2 M1")
}

func (S2) M2() {
	fmt.Println("S2 M2")
}

func showStructEmbedInterface2() {
	t := T2{
		Interface2: S2{},
	}
	t.M1()
	t.M2()
}

func main() {
	showStructEmbedInterface2()
}
