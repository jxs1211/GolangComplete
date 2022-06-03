package main

import "fmt"

func foo(a, b int) (x, y int) {
	defer func() {
		x = x * 5
		y = y * 10
	}()

	x = a + 5
	y = b + 6
	return
}

func showModifyReturnVarInDefer() {
	x, y := foo(1, 2)
	fmt.Println("x=", x, "y=", y)
}

func foo2(a, b int) (x, y int) {
	defer func() {
		x *= 10
		y *= 20
	}()
	a += 1
	b += 2
	x, y = a, b
	return
}

func showModifyReturnVarInDefer2() {
	x, y := foo2(1, 2)
	fmt.Println(x, y)
}

func main() {
	showModifyReturnVarInDefer2()
}
