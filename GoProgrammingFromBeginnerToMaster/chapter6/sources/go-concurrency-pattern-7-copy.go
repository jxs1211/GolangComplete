package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Shutdowner2 interface {
	Shutdown(timeout time.Duration) error
}

type ShutdownFunc func(timeout time.Duration) error

func (f ShutdownFunc) Shutdown(timeout time.Duration) error {
	return f(timeout)
}

func fork(timeout time.Duration, shutdowners ...Shutdowner2) chan struct{} {
	done := make(chan struct{})
	var wg sync.WaitGroup
	for i, s := range shutdowners {
		wg.Add(1)
		go func(id int) {
			defer func() {
				wg.Done()
				fmt.Printf("worker-%d done\n", id)
			}()
			s.Shutdown(timeout)
		}(i)
	}

	go func() {
		wg.Wait()
		done <- struct{}{}
	}()
	return done
}

func ConcurrentShutdown2(timeout time.Duration, shutdowners ...Shutdowner2) error {
	done := fork(timeout, shutdowners...)

	t := time.NewTimer(5 * time.Second)
	defer t.Stop()
	select {
	case <-done:
		fmt.Println("all worker done")
		return nil
	case <-t.C:
		fmt.Println("timeout")
		return errors.New("timeout")
	}
}

func SequentialShutdown2(timeout time.Duration, shutdowners ...Shutdowner2) error {
	done := make(chan struct{})
	start := time.Now()
	// elapsed := time.Now() - start
	var left time.Duration

	t := time.NewTimer(timeout)
	for i, s := range shutdowners {
		elapsed := time.Since(start) // calculate previous goroutine's time duration
		left = timeout - elapsed
		go func(id int) {
			begin := time.Now()
			defer func() {
				fmt.Printf("worker-%d: done, cost: %ds\n", id, time.Since(begin)/time.Second)
				done <- struct{}{}
			}()
			s.Shutdown(left)
		}(i)
		t.Reset(left)
		select {
		case <-t.C:
			return errors.New("timeout")
		case <-done:
			// continue to exec next worker
		}
	}
	return nil
}
