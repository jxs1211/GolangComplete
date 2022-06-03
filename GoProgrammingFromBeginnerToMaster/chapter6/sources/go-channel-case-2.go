package main

import (
	"fmt"
	"time"
)

type signal struct{}

func worker() {
	println("worker is working...")
	time.Sleep(1 * time.Second)
}

func spawn(f func()) <-chan signal {
	c := make(chan signal)
	go func() {
		println("worker start to work...")
		f()
		c <- signal(struct{}{})
	}()
	return c
}

func worker2() {
	time.Sleep(1 * time.Second)
	fmt.Println("call worker2")
}

func spawn2(f func()) <-chan struct{} {
	c := make(chan struct{})
	go func() {
		f()
		c <- struct{}{}
	}()
	return c
}

func showGoChanCase2() {
	c := spawn2(worker2)
	v := <-c
	fmt.Println(v)
}

func main() {
	showGoChanCase2()
}
