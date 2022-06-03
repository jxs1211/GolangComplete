package main

import (
	"fmt"
)

type Counter22 struct {
	c chan int
	i int
}

var semaphore = 10
var counter22 Counter22

func addCount() int {
	return <-counter22.c
}

func InitCounter2() <-chan struct{} {
	signal := make(chan struct{})
	counter22 = Counter22{
		c: make(chan int),
	}
	go func() {
		for {
			if counter22.i >= semaphore {
				signal <- struct{}{}
				fmt.Println("counter done")
				return
			}
			counter22.i++
			counter22.c <- counter22.i
		}
	}()
	fmt.Println("counter init ok")
	return signal
}

func init() {
	// InitCounter2()
}

func main() {
	signal := InitCounter2()
	for i := 0; i < semaphore; i++ {
		go func(id int) {
			v := addCount()
			fmt.Printf("goroutine %d: increase done: %d\n", id, v)
		}(i)
	}
	// time.Sleep(time.Second * 5)
	<-signal
	fmt.Println("main goroutine done")

}
