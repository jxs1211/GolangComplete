package main

import (
	"errors"
	"fmt"
	"sync"
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

func spawnGroup(n int, f func(args ...interface{}), args ...interface{}) chan struct{} {
	c := make(chan struct{})
	var wg sync.WaitGroup

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			name := fmt.Sprintf("worker-%d:", i)
			f(args...)
			println(name, "done")
			wg.Done() // worker done!
		}(i)
	}

	go func() {
		wg.Wait()
		c <- struct{}{}
	}()

	return c
}

func showConcurrencyPattern3() {
	done := spawnGroup(5, worker, 1)
	<-done
	fmt.Println("main goroutine done")
}

func bar31(args ...interface{}) error {
	if len(args) == 0 {
		return errors.New("invalid args")
	}
	v, ok := args[0].(int)
	if !ok {
		return errors.New("invalid int args")
	}
	time.Sleep(time.Second * time.Duration(v))
	fmt.Printf("call bar elasped: %ds\n", v)
	return nil
}

func spawn31(nums int, f func(args ...interface{}) error, args ...interface{}) <-chan struct{} {
	done := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < nums; i++ {
		wg.Add(1)
		go func(id int) {
			fmt.Printf("goroutinue %d is working\n", id)
			f(args...)
			wg.Done()
		}(i)
	}

	// monitor worker goroutine
	go func() {
		wg.Wait()
		done <- struct{}{}
	}()
	return done
}

func showConcurrencyPattern31() {
	done := spawn31(5, bar31, 1)
	<-done
	fmt.Println("main goroutine done")
}

func main() {
	showConcurrencyPattern3()
	// showConcurrencyPattern31()
}
