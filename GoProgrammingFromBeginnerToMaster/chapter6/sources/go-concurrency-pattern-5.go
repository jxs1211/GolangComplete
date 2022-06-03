package main

import (
	"fmt"
	"time"
)

func worker(j int) {
	fmt.Println("worker start")
	time.Sleep(time.Second * (time.Duration(j)))
	fmt.Println("worker done")
}

func spawn(f func(int)) chan string {
	quit := make(chan string)
	go func() {
		defer fmt.Println("worker exit")
		var job chan int // 模拟job channel
	loopout:
		for {
			select {
			case j := <-job:
				// default:
				f(j) // simulate worker running
			case <-quit:
				quit <- "ok"
				break loopout
			}
		}
	}()
	return quit
}

func main() {
	quit := spawn(worker)
	println("spawn a worker goroutine")

	time.Sleep(5 * time.Second)

	// notify the child goroutine to exit
	println("notify the worker to exit...")
	quit <- "exit"

	timer := time.NewTimer(time.Second * 5)
	defer timer.Stop()
	select {
	case status := <-quit:
		println("worker done:", status)
	case <-timer.C:
		println("wait worker exit timeout")
	}
}
