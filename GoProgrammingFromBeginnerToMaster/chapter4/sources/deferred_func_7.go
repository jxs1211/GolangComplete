package main

import "fmt"

func foo1() {
	for i := 0; i <= 3; i++ {
		defer fmt.Println(i)
	}
}

func foo2() {
	for i := 0; i <= 3; i++ {
		defer func(n int) {
			fmt.Println(n)
		}(i)
	}
}

func foo3() {
	for i := 0; i <= 3; i++ {
		defer func() {
			fmt.Println(i)
		}()
	}
}

func foo11() {
	for i := 0; i <= 3; i++ {
		defer fmt.Println(i)
	}
}

func foo12() {
	for i := 0; i <= 3; i++ {
		defer func(a int) {
			fmt.Println(a)
		}(i)
	}
}

func foo13() {
	for i := 0; i <= 3; i++ {
		defer func() {
			fmt.Println(i)
		}()
	}
}

func showWhenEvaluationWhenDeferRegistering() {
	fmt.Println("foo1 result:")
	foo11()
	fmt.Println("foo2 result:")
	foo12()
	fmt.Println("foo3 result:")
	foo13()
}

func main() {
	showWhenEvaluationWhenDeferRegistering()
}
