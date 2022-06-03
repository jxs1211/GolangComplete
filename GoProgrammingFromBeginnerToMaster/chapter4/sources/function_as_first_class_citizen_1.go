package main

import "fmt"

type binaryCalcFunc func(int, int) int

func showFunc() {
	var i interface{} = binaryCalcFunc(func(x, y int) int { return x + y })
	c := make(chan func(int, int) int, 10)
	fns := []binaryCalcFunc{
		func(x, y int) int { return x + y },
		func(x, y int) int { return x - y },
		func(x, y int) int { return x * y },
		func(x, y int) int { return x / y },
		func(x, y int) int { return x % y },
	}

	c <- func(x, y int) int {
		return x * y
	}

	fmt.Println(fns[0](5, 6))
	f := <-c
	fmt.Println(f(7, 10))
	v, ok := i.(binaryCalcFunc)
	if !ok {
		fmt.Println("type assertion error")
		return
	}

	fmt.Println(v(17, 7))
}

type Calculate func(x, y int) int

func showFunc2() {
	var i interface{} = Calculate(func(x, y int) int { return x + y })
	// i.(Calculate)(10, 10)
	fns := []Calculate{
		func(x, y int) int { return x + y },
		func(x, y int) int { return x - y },
		func(x, y int) int { return x * y },
		func(x, y int) int { return x / y },
		func(x, y int) int { return x % y },
	}
	fmt.Println(fns[0](1, 10))
	ch := make(chan Calculate, 10)
	ch <- fns[1]
	ch <- func(x, y int) int { return x / y }
	v := <-ch
	fmt.Println(v(10, 1))
	v = <-ch
	fmt.Println(v(9, 3))
	ch2 := make(chan func(int, int) int, 10)
	ch2 <- fns[4]
	v2 := <-ch2
	fmt.Println(v2(5, 3))
	v, ok := i.(Calculate)
	if !ok {
		fmt.Println("type assertion error")
	}

	fmt.Println(v(10, 10))
}

func main() {
	showFunc2()
}
