package main

import "fmt"

func times(x, y int) int {
	return x * y
}

func partialTimes(x int) func(int) int {
	return func(y int) int {
		return times(x, y)
	}
}

func showFuncCurring() {
	timesTwo := partialTimes(2)
	timesThree := partialTimes(3)
	timesFour := partialTimes(4)
	fmt.Println(timesTwo(5))
	fmt.Println(timesThree(5))
	fmt.Println(timesFour(5))
}

func partialTimes2(x int) func(int) int {
	return func(y int) int {
		return x * y
	}
}

func showFuncCurring2() {
	times2 := partialTimes2(2)
	times3 := partialTimes2(3)
	times4 := partialTimes2(4)
	fmt.Println(times2(2))
	fmt.Println(times3(2))
	fmt.Println(times4(2))
}

func main() {
	showFuncCurring2()
}
