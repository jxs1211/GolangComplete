package main

import (
	"fmt"
	"sync"
	"time"
)

type signal struct{}

func worker3(i int) {
	fmt.Printf("worker %d: is working...\n", i)
	time.Sleep(1 * time.Second)
	fmt.Printf("worker %d: works done\n", i)
}

func spawnGroup(f func(i int), num int, groupSignal <-chan signal) <-chan signal {
	c := make(chan signal)
	var wg sync.WaitGroup

	for i := 0; i < num; i++ {
		wg.Add(1)
		go func(i int) {
			<-groupSignal
			fmt.Printf("worker %d: start to work...\n", i)
			f(i)
			wg.Done()
		}(i + 1)
	}

	go func() {
		wg.Wait()
		c <- signal(struct{}{})
	}()
	return c
}

func showGoChanCase3() {
	fmt.Println("start a group of workers...")
	groupSignal := make(chan signal)
	c := spawnGroup(worker3, 5, groupSignal)
	time.Sleep(5 * time.Second)
	fmt.Println("the group of workers start to work...")
	close(groupSignal)
	<-c
	fmt.Println("the group of workers work done!")
}

func worker33(id int) {
	fmt.Printf("worker %d: is working...\n", id)
	time.Sleep(1 * time.Second)
	fmt.Printf("worker %d: works done\n", id)
}

func spawnGroup2(f func(i int), nums int, workerSignal <-chan struct{}) <-chan struct{} {
	c := make(chan struct{})
	// monWorkerSignal := make(chan struct{}, nums)
	// monwSignal := make(chan struct{})
	var wg sync.WaitGroup

	for i := 0; i < nums; i++ {
		wg.Add(1)
		go func(id int) {
			<-workerSignal
			f(id)
			wg.Done()
			// monWorkerSignal <- struct{}{}
			// // monwSignal <- struct{}{}
		}(i + 1)
	}

	go func() {
		wg.Wait()
		fmt.Println("monitor goroutinue done")
		c <- struct{}{}
	}()
	return c
}

func showGoChanCase33() {
	workerSignal := make(chan struct{})
	monSignal := spawnGroup2(worker33, 5, workerSignal)
	time.Sleep(5 * time.Second)
	close(workerSignal)
	<-monSignal
	fmt.Println("all task done")
}

func main() {
	showGoChanCase33()
}

// notify: main goroutine -> monitor goroutine
//                        -> worker goroutine

// response: main goroutine <- monitor goroutine <- worker goroutine
