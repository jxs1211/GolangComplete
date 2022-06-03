package main

import (
	"fmt"
	"sync"
	"time"
)

type signal struct{}

var ready bool

func worker(i int) {
	fmt.Printf("worker %d: is working...\n", i)
	time.Sleep(1 * time.Second)
	fmt.Printf("worker %d: works done\n", i)
}

func spawnGroup(f func(i int), num int, mu *sync.Mutex) <-chan signal {
	c := make(chan signal)
	var wg sync.WaitGroup

	for i := 0; i < num; i++ {
		wg.Add(1)
		go func(i int) {
			for {
				mu.Lock()
				if !ready {
					mu.Unlock()
					time.Sleep(100 * time.Millisecond)
					continue
				}
				mu.Unlock()
				fmt.Printf("worker %d: start to work...\n", i)
				f(i)
				wg.Done()
				return
			}
		}(i + 1)
	}

	go func() {
		wg.Wait()
		c <- signal(struct{}{})
	}()
	return c
}

func showSyncPkg1() {
	fmt.Println("start a group of workers...")
	mu := &sync.Mutex{}
	c := spawnGroup(worker, 5, mu)

	time.Sleep(5 * time.Second) // 模拟ready前的准备工作
	fmt.Println("the group of workers start to work...")

	mu.Lock()
	ready = true
	mu.Unlock()

	<-c
	fmt.Println("the group of workers work done!")

}

func work(duration int) {
	fmt.Println("worker start")
	time.Sleep(time.Millisecond * time.Duration(duration))
	fmt.Println("worker end")
}

var ready2 bool
var mutex sync.Mutex

func spawnGoroutine(f func(i int), nums int, cond *sync.Cond) <-chan struct{} {
	done := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < nums; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cond.L.Lock()
			if !ready2 {
				cond.Wait()
			}
			cond.L.Unlock()
			f(i)
			// for {
			// 	mu.Lock()
			// 	if !ready2 {
			// 		mu.Unlock()
			// 		// fmt.Printf("goroutine-%d: not ready..", i)
			// 		time.Sleep(time.Second)
			// 		continue
			// 	}
			// 	mu.Unlock()
			// 	f(i)
			// 	return
			// }
		}(i)
	}

	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	return done
}

func showSyncPkg2() {
	// mu := &sync.Mutex{}
	cond := sync.NewCond(&sync.Mutex{})
	c := spawnGoroutine(work, 5, cond)
	fmt.Println("exec group workers")
	time.Sleep(5 * time.Second)

	cond.L.Lock()
	ready2 = true
	cond.Broadcast()
	fmt.Println("trigger worker")
	cond.L.Unlock()

	<-c
	fmt.Println("all worker done")
}

func main() {
	// showSyncPkg1()
	showSyncPkg2()
}
