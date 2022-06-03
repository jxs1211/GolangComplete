package main

import (
	"errors"
	"fmt"
	"time"
)

var OK = errors.New("ok")

func worker(args ...interface{}) error {
	if len(args) == 0 {
		return errors.New("invalid args")
	}
	interval, ok := args[0].(int)
	if !ok {
		return errors.New("invalid interval arg")
	}

	time.Sleep(time.Second * (time.Duration(interval)))
	return OK
}

func spawn(f func(args ...interface{}) error, args ...interface{}) chan error {
	c := make(chan error)
	go func() {
		c <- f(args...)
	}()
	return c
}

func showConcurrencyPattern2() {
	done := spawn(worker, 5)
	println("spawn worker1")
	err := <-done
	fmt.Println("worker1 done:", err)
	done = spawn(worker)
	println("spawn worker2")
	err = <-done
	fmt.Println("worker2 done:", err)
}

func bar21(args ...interface{}) error {
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

func spawn21(f func(args ...interface{}) error, args ...interface{}) <-chan error {
	done := make(chan error)
	go func() {
		err := f(args...)
		done <- err
	}()
	return done
}

func showConcurrencyPattern21() {
	done := spawn21(bar21, 3)
	err := <-done
	if err != nil {
		fmt.Println(err)
	}
	done = spawn21(bar21)
	err = <-done
	if err != nil {
		fmt.Println(err)
	}
	done = spawn21(bar21, "2")
	err = <-done
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("main goroutine done")
}

func main() {
	// showConcurrencyPattern2()
	showConcurrencyPattern21()
}
