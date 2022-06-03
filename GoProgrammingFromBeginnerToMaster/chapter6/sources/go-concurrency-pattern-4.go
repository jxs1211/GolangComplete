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

func showConcurrencyPattern4() {
	done := spawnGroup(5, worker, 30)
	println("spawn a group of workers")

	timer := time.NewTimer(time.Second * 5)
	defer timer.Stop()
	select {
	case <-timer.C:
		println("wait group workers exit timeout!")
	case <-done:
		println("group workers done")
	}
}

func bar41(id int, args ...interface{}) error {
	if len(args) == 0 {
		return errors.New("invalid args")
	}
	v, ok := args[0].(int)
	if !ok {
		return errors.New("invalid int args")
	}
	fmt.Printf("goroutinue %d is working\n", id)
	time.Sleep(time.Second * time.Duration(v))
	fmt.Printf("\tcall bar elasped: %ds\n", v)
	fmt.Printf("goroutinue %d is done\n", id)
	return nil
}

func bar42(id int, args ...interface{}) (int, error) {
	if len(args) == 0 {
		return 0, errors.New("invalid args")
	}
	v, ok := args[0].(int)
	if !ok {
		return 0, errors.New("invalid int args")
	}
	if id%2 != 0 {
		return 0, errors.New("goroutine exec error")
	}
	return id * v, nil
}

func spawn41(nums int, f func(id int, args ...interface{}) error, args ...interface{}) <-chan struct{} {
	done := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < nums; i++ {
		wg.Add(1)
		go func(id int) {
			f(id, args...)
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

func showConcurrencyPattern41() {
	done := spawn41(5, bar41, 10)
	t := time.NewTimer(5 * time.Second)
	select {
	case <-done:
		fmt.Println("worker all done")
	case <-t.C:
		fmt.Println("worker timeout")
	}
	fmt.Println("main goroutine done")
}

type result42 struct {
	r   int
	err error
}

func spawn42(nums int, f func(id int, args ...interface{}) (int, error), args ...interface{}) <-chan result42 {
	done := make(chan result42)

	// var wg sync.WaitGroup
	for i := 0; i < nums; i++ {
		// wg.Add(1)
		go func(id int) {
			r, err := f(id, args...)

			done <- result42{r: r, err: err}
			// wg.Done()
		}(i)
	}

	return done
}

func showConcurrencyPattern42() {
	nums := 5
	done := spawn42(nums, bar42, 2)
	res := make([]result42, 0, nums)
	for i := 0; i < nums; i++ {
		r := <-done
		res = append(res, r)
	}
	fmt.Println(res)
	fmt.Println("main goroutine done")
}

func main() {
	// showConcurrencyPattern41()
	showConcurrencyPattern42()
}
