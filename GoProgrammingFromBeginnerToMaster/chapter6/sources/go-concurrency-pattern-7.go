package main

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type GracefullyShutdowner interface {
	Shutdown(waitTimeout time.Duration) error
}

type ShutdownerFunc func(time.Duration) error

func (f ShutdownerFunc) Shutdown(waitTimeout time.Duration) error {
	return f(waitTimeout)
}

func ConcurrentShutdown(waitTimeout time.Duration, shutdowners ...GracefullyShutdowner) error {
	c := make(chan struct{})

	go func() {
		var wg sync.WaitGroup
		for _, g := range shutdowners {
			wg.Add(1)
			go func(shutdowner GracefullyShutdowner) {
				defer wg.Done()
				shutdowner.Shutdown(waitTimeout)
			}(g)
		}
		wg.Wait()
		c <- struct{}{}
	}()

	timer := time.NewTimer(waitTimeout)
	defer timer.Stop()

	select {
	case <-c:
		return nil
	case <-timer.C:
		return errors.New("wait timeout")
	}
}

func SequentialShutdown(waitTimeout time.Duration, shutdowners ...GracefullyShutdowner) error {
	start := time.Now()
	var left time.Duration
	timer := time.NewTimer(waitTimeout)

	for i, g := range shutdowners {
		elapsed := time.Since(start)
		left = waitTimeout - elapsed
		fmt.Printf("id: %d, left: %ds\n", i, left/time.Second)
		c := make(chan struct{})
		go func(id int, shutdowner GracefullyShutdowner) {
			begin := time.Now()
			fmt.Printf("goroutine %d start, begin: %v\n", id, begin)
			defer fmt.Printf("goroutine %d end\n", id)
			shutdowner.Shutdown(left)
			c <- struct{}{}
			fmt.Printf("goroutine %d cost: %ds\n", id, time.Since(begin)/time.Second)
		}(i, g)

		timer.Reset(left)
		select {
		case <-c:
			//continue
		case <-timer.C:
			fmt.Println("timeout!")
			return errors.New("wait timeout")
		}
	}
	return nil
}


