package main

import "fmt"

type BinaryAdder interface {
	Add(int, int) int
}

type MyAdderFunc func(int, int) int

func (f MyAdderFunc) Add(x, y int) int {
	return f(x, y)
}

func MyAdd(x, y int) int {
	return x + y
}

func showFuncAsType() {
	var i BinaryAdder = MyAdderFunc(MyAdd)
	fmt.Println(i.Add(5, 6))
}

type Adder interface {
	Add(int, int) int
}

type AdderFunc func(x, y int) int

func (a AdderFunc) Add(x, y int) int {
	return a(x, y)
}

func doAdd(x, y int) int {
	return x + y
}

func showFuncAsType2() {
	var adder Adder = AdderFunc(doAdd)
	fmt.Println(adder.Add(1, 11))
}

func main() {
	showFuncAsType()
}
