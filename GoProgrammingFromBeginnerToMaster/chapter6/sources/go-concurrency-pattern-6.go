package main

import (
	"fmt"
	"sync"
	"time"
)

func worker6(j int) {
	time.Sleep(time.Second * (time.Duration(j)))
	fmt.Println("exec...")
}

func spawnGroup6(n int, f func(int)) chan struct{} {
	quit := make(chan struct{})
	job := make(chan int)
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer func() {
				wg.Done()
				fmt.Printf("worker-%d done\n", i)
			}() // 保证wg.Done在goroutine退出前被执行
			fmt.Printf("worker-%d start\n", i)
			for {
				j, ok := <-job
				if !ok {
					// println(name, "done")
					return
				}
				// do the job
				worker6(j)
			}
		}(i)
	}

	go func() {
		<-quit
		close(job) // 广播给所有新goroutine
		wg.Wait()
		quit <- struct{}{}
	}()

	return quit
}

func showConcurrencyPattern6() {
	quit := spawnGroup6(5, worker6)
	println("spawn a group of workers")

	time.Sleep(3 * time.Second)
	// notify the worker goroutine group to exit
	println("notify the worker group to exit...")
	quit <- struct{}{}

	timer := time.NewTimer(time.Second * 3)
	defer timer.Stop()
	select {
	case <-timer.C:
		println("wait group workers exit timeout!")
	case <-quit:
		println("group workers done")
	}
}

func bar61(duration int) {
	fmt.Println("bar61")
	time.Sleep(time.Second * time.Duration(duration))
	fmt.Println("bar61 end")
}

func spawn61(nums int, f func(d int)) chan struct{} {
	done := make(chan struct{})
	job := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < nums; i++ {
		wg.Add(1)
		go func(id int) {
			defer func() {
				wg.Done()
				fmt.Printf("worker-%d done\n", id)
			}()
			fmt.Printf("worker-%d start\n", id)
			for {
				_, ok := <-job
				if !ok {
					return
				}
				f(7)
				return
			}
		}(i)
	}

	go func() {
		defer fmt.Println("mon goroutine done")
		<-done
		for i := 0; i < nums; i++ {
			job <- struct{}{}
		}
		close(job)
		wg.Wait()
		done <- struct{}{}
	}()
	return done
}

func showConcurrencyPattern61() {
	done := spawn61(5, bar61)
	time.Sleep(3 * time.Second)
	fmt.Println("notify")
	done <- struct{}{}

	t := time.NewTimer(3 * time.Second)
	defer t.Stop()
	select {
	case <-done:
		fmt.Println("workers all done")
	case <-t.C:
		fmt.Println("timeout")
	}
}

func main() {
	// showConcurrencyPattern6()
	showConcurrencyPattern61()
}
