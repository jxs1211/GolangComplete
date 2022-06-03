package main

import (
	"fmt"
	"sync"
	"time"
)

type signal struct{}

func worker(i int, quit <-chan signal) {
	fmt.Printf("worker %d: is working...\n", i)
LOOP:
	for {
		select {
		default:
			// 模拟worker工作
			time.Sleep(1 * time.Second)

		case <-quit:
			break LOOP
		}
	}
	fmt.Printf("worker %d: works done\n", i)
}

func spawnGroup4(f func(int, <-chan signal), num int, groupSignal <-chan signal) <-chan signal {
	c := make(chan signal)
	var wg sync.WaitGroup

	for i := 0; i < num; i++ {
		wg.Add(1)
		go func(i int) {
			fmt.Printf("worker %d: start to work...\n", i)
			f(i, groupSignal)
			wg.Done()
		}(i + 1)
	}

	go func() {
		wg.Wait()
		c <- signal(struct{}{})
	}()
	return c
}

func worker44(id int, signal <-chan struct{}) {
	fmt.Printf("worker %d: is working...\n", id)
loopout:
	for {
		select {
		default:
			time.Sleep(1 * time.Second)
		case <-signal:
			break loopout
		}
	}
	fmt.Printf("worker %d: works done\n", id)
}

func spawnGroup44(f func(i int, c <-chan struct{}), nums int, workerSignal <-chan struct{}) <-chan struct{} {
	c := make(chan struct{})
	var wg sync.WaitGroup

	for i := 0; i < nums; i++ {
		wg.Add(1)
		go func(id int) {
			// <-workerSignal
			f(id, workerSignal)
			wg.Done()
		}(i + 1)
	}

	go func() {
		wg.Wait()
		fmt.Println("monitor goroutinue done")
		c <- struct{}{}
	}()
	return c
}

func showGoChanCase44() {
	workerSignal := make(chan struct{})
	monSignal := spawnGroup44(worker44, 5, workerSignal)
	time.Sleep(5 * time.Second)
	fmt.Println("notify the group of workers to exit...")
	close(workerSignal)
	<-monSignal
	fmt.Println("all task done")
}

func showGoChanCase4() {
	fmt.Println("start a group of workers...")
	groupSignal := make(chan signal)
	c := spawnGroup4(worker, 5, groupSignal)
	fmt.Println("the group of workers start to work...")

	time.Sleep(5 * time.Second)
	// 通知workers退出
	fmt.Println("notify the group of workers to exit...")
	close(groupSignal)
	<-c
	fmt.Println("the group of workers work done!")
}

func main() {
	// showGoChanCase4()
	showGoChanCase44()
}
