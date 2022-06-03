package main

import "fmt"

var c = make(chan int)
var a string

func f() {
	a = "hello, world"
	<-c
}

func showGoChanCase1() {
	go f()
	c <- 5
	println(a)
}

func bar() <-chan int {
	ch := make(chan int)
	go func() {
		fmt.Println("hello bar")
		ch <- 1
	}()
	return ch
}

func showGoChanCase11() {
	ch := bar()
	v, ok := <-ch
	if !ok {
		fmt.Println("nothing received")
		return
	}
	fmt.Println("receive: ", v)

}

func main() {
	// showGoChanCase1()
	showGoChanCase11()
}
