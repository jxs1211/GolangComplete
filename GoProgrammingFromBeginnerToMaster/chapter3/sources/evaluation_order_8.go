package main

import (
	"fmt"
	"time"
)

func getAReadOnlyChannel() <-chan int {
	fmt.Println("invoke getAReadOnlyChannel")
	c := make(chan int)

	go func() {
		time.Sleep(3 * time.Second)
		c <- 1
	}()

	return c
}

func getASlice() *[5]int {
	fmt.Println("invoke getASlice")
	var a [5]int
	return &a
}

func getAWriteOnlyChannel() chan<- int {
	fmt.Println("invoke getAWriteOnlyChannel")
	return make(chan int)
}

func getANumToChannel() int {
	fmt.Println("invoke getANumToChannel")
	return 2
}

func showChan() {
	select {
	//recv from channel
	case (getASlice())[0] = <-getAReadOnlyChannel():
		fmt.Println("recv something from a readonly channel")

	//send to channel
	case getAWriteOnlyChannel() <- getANumToChannel():
		fmt.Println("send something to a writeonly channel")
	}
}

func getSlice() []int {
	fmt.Println("getSlice")
	s := make([]int, 1)
	return s
}

func getReadChan() <-chan int {
	fmt.Println("getReadChan")
	ch := make(chan int)
	go func() {
		for {
			time.Sleep(3 * time.Second)
			ch <- 1
		}
	}()
	return ch
}

func getWriteChan() chan<- int {
	fmt.Println("getWriteChan")
	return make(chan int)
}

func sendAValue() int {
	fmt.Println("sendAValue")
	return 1
}

func showChan2() {
	select {
	case getSlice()[0] = <-getReadChan():
		fmt.Println("get element from a read channel")
	case getWriteChan() <- sendAValue():
		fmt.Println("send element to a write channel")
	}
}

func main() {
	showChan2()
}
