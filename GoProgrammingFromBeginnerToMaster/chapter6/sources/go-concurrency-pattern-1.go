package main

import (
	"fmt"
	"time"
)

func worker(args ...interface{}) {
	if len(args) == 0 {
		return
	}
	interval, ok := args[0].(int)
	if !ok {
		return
	}

	time.Sleep(time.Second * (time.Duration(interval)))
}

func spawn(f func(args ...interface{}), args ...interface{}) chan struct{} {
	c := make(chan struct{})
	go func() {
		f(args...)
		c <- struct{}{}
	}()
	return c
}

func showConcurrencyPattern1() {
	done := spawn(worker, 5)
	println("spawn a worker goroutine")
	<-done
	println("worker done")
}

func bar11(args ...interface{}) {
	v, ok := args[0].(int)
	if !ok {
		return
	}
	time.Sleep(time.Second * time.Duration(v))
	fmt.Println("call bar elasped: ", v)

}

func spawn11(f func(args ...interface{}), args ...interface{}) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		f(args...)
		done <- struct{}{}
	}()
	return done
}

func showConcurrencyPattern11() {
	done := spawn11(bar11, 2)
	<-done
	fmt.Println("main goroutine done")
}

func main() {
	showConcurrencyPattern11()
}
