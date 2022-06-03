package main

import "fmt"

func sum(args ...int) int {
	var total int

	for _, v := range args {
		total += v
	}

	return total
}

func showVariadic() {
	a, b, c := 1, 2, 3
	println(sum(a, b, c))
	nums := []int{4, 5, 6}
	println(sum(nums...))
}

func multiply(ints ...int) int { // function layer
	total := 1
	for _, v := range ints { // inner layer
		total *= v
	}
	return total
}

func main() {
	fmt.Println(multiply(1, 2, 3)) // outlayer
	fmt.Println(multiply([]int{1, 2, 3}...))
}
