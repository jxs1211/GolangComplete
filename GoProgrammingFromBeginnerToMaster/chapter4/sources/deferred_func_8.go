package main

import "fmt"

func foo1() {
	sl := []int{1, 2, 3}
	defer func(a []int) {
		fmt.Println(a)
	}(sl)

	sl = []int{3, 2, 1}
	_ = sl
}
func foo2() {
	sl := []int{1, 2, 3}
	defer func(p *[]int) {
		fmt.Println(*p)
	}(&sl)

	sl = []int{3, 2, 1}
	_ = sl
}

func foo11() {
	sl := []int{1, 2, 3}
	defer func(s []int) {
		fmt.Println(s)
	}(sl)
	sl = []int{3, 2, 1}
}

func foo12() {
	sl := []int{1, 2, 3}
	defer func(s *[]int) {
		fmt.Println(*s)
	}(&sl)
	sl = []int{3, 2, 1}
}

func showWhenEvaluationWhenDeferRegistering() {
	foo11()
	foo12()
}

func main() {
	showWhenEvaluationWhenDeferRegistering()
}
