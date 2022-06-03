package main

import (
	"fmt"
	"math"
)

func Foo() {
	{ //隐式代码块1 开始
		type bar struct{} // 类型标识符bar作用域 开始
		{                 //隐式代码块2 开始
			a := 5 // 变量标识符a作用域 开始
			{
				fmt.Println(a)
			}
			// 变量标识符a作用域 结束
		} //隐式代码块2 结束
		// 类型标识符bar作用域 结束
	} //隐式代码块1 结束
}

func showIf() {
	if a := 1; true {
		fmt.Println(a)
	}
}

func showIfEqual() {
	{
		a := 1
		if true {
			fmt.Println(a)
		}
	}
}

func showIfElse() {
	if a := 1; true {
		fmt.Println(a)
	} else {
		fmt.Println(math.MaxInt)
	}
}

func showIfElseEqual() {
	{
		a := 1
		if true {
			fmt.Println(a)
		} else {
			fmt.Println(math.MaxInt)
		}
	}
}

func showIfElseIfElse() {
	if a := 1; false {
	} else if b := 2; false {
	} else if c := 3; false {
	} else {
		fmt.Println(a, b, c)
	}
}

func showIfElseIfElseEqual() {
	{
		a := 1
		if false {
		} else {
			{
				b := 2
				if false {
				} else {
					{
						c := 3
						if false {
						} else {
							fmt.Println(a, b, c)
						}
					}
				}
			}
		}
	}
}

func showFor1() {
	for i := 1; i < 3; i++ {
		fmt.Println(i)
	}
}

func showFor1Equal() {
	{
		i := 1
		for ; i < 3; i++ {
			fmt.Println(i)
		}
	}
}

func showFor2() {
	s := []int{1, 2, 3}
	for i, v := range s {
		fmt.Println(i, v)
	}
}

func showFor2Equal() {
	s := []int{1, 2, 3}
	{
		i, v := 0, 0
		for i, v = range s {
			fmt.Println(i, v)
		}
	}
}

func showSwitch1() {
	switch x, y := 1, 1; x + y {
	case 1:
		a := 10
		fmt.Println("a: ", a)
		fmt.Println("res: ", 1)
		fallthrough
	case 2:
		a := 20
		fmt.Println("a: ", a)
		fmt.Println("res: ", 2)
		fallthrough
	default:
		a := 30
		fmt.Println("a: ", a)
		fmt.Println("no result")
	}
}

func showSwitch1Equal() {
	{
		x, y := 1, 1
		switch x + y {
		case 1:
			{
				a := 10
				fmt.Println("a: ", a)
				fmt.Println("res: ", 1)
			}
			fallthrough
		case 2:
			{
				a := 20
				fmt.Println("a: ", a)
				fmt.Println("res: ", 2)
			}
			fallthrough
		default:
			{
				a := 30
				fmt.Println("a: ", a)
				fmt.Println("no result")
			}
		}
	}
}

func showSelect() {
	c1 := make(chan int)
	c2 := make(chan int, 1)
	c2 <- 1

	select {
	case c1 <- 1:
		fmt.Println("send to channel")
	case a := <-c2:
		_ = a
		fmt.Println("receive from channel")
	default:
		fmt.Println("default")
	}
}

func showSelectEqual() {
	c1 := make(chan int)
	c2 := make(chan int, 1)
	c2 <- 1

	select {
	case c1 <- 1:
		{
			fmt.Println("send to channel")
		}
	case a := <-c2: //如果该条件被选中
		{
			// a := <-c2
			_ = a
			fmt.Println("receive from channel")
		}
	default:
		{
			fmt.Println("default")
		}
	}
}

func main() {
	showSelectEqual()
	showSelect()
	// showSwitch1Equal()
	// showSwitch1()
	// showFor2Equal()
	// showFor2()
	// showFor1()
	// showFor1Equal()
	// showIfElseIfElse()
	// showIfElseIfElseEqual()
	// showIfElse()
	// showIfElseEqual()
	// Foo()
	// showIf()
	// showIfEqual()
}
