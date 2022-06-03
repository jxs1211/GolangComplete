package main

import (
	"fmt"
	"time"
)

func demo1() {
	var m = [...]int{1, 2, 3, 4, 5}

	for i, v := range m {
		go func() {
			time.Sleep(time.Second * 3)
			fmt.Println(i, v)
		}()
	}

	time.Sleep(time.Second * 10)
}

func demo2() {
	var m = [...]int{1, 2, 3, 4, 5}

	for i, v := range m {
		go func(i, v int) {
			time.Sleep(time.Second * 3)
			fmt.Println(i, v)
		}(i, v)
	}

	time.Sleep(time.Second * 10)
}

func showForLoopLocalVariableReuse() {
	s := []int{1, 2, 3}
	for i, v := range s {
		go func() {
			fmt.Println(i, v)
		}()
	}
	time.Sleep(time.Second)
}

func showForLoopLocalVariableReuse2() {
	s := []int{1, 2, 3}
	for i, v := range s {
		go func(i, v int) {
			fmt.Println(i, v)
		}(i, v)
	}
	time.Sleep(time.Second)
}

func main() {
	showForLoopLocalVariableReuse()
	showForLoopLocalVariableReuse2()
	// demo1()
	// demo2()
}
