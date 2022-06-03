package main

import (
	"fmt"
	"sync"
	"time"
)

func bar() {
	fmt.Println("raise a panic")
	panic(-1)
}

func foo() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("recovered from a panic")
		}
	}()
	bar()
}

func showRecoverPanic() {
	foo()
	fmt.Println("main exit normally")
}

func bar2() {
	fmt.Println("bar2")
	panic(-1)
}

func foo2() {
	defer func() {
		if e := recover(); e != nil {
			fmt.Println("recover from panic: ", e)
		}
	}()
	bar2()
}

func showRecoverPanic2() {
	foo2()
	fmt.Println("main exit")

}

var mu sync.Mutex

func bizOp() {
	fmt.Println("do business")
	panic("ops! raise a panic here")

}

func BehaveWithoutDefer() {
	fmt.Println("exec BehaveWithoutDefer")
	mu.Lock()
	bizOp() // mutex can't be released if panic here
	mu.Unlock()
}

func BehaveWithDefer() {
	fmt.Println("exec BehaveWithDefer")
	mu.Lock()
	defer mu.Unlock()
	bizOp()
}

func Println() {
	time.Sleep(time.Second)
	fmt.Println("just print....")
}

func main() {
	go Println()
	go BehaveWithDefer()
	time.Sleep(time.Second * 3)
	// showRecoverPanic2()
}
