package main

import (
	"fmt"
	"time"
)

func doIteration(m map[int]int) {
	for k, v := range m {
		_ = fmt.Sprintf("[%d, %d] ", k, v)
	}
}
func doWrite(m map[int]int) {
	for k, v := range m {
		m[k] = v + 1
	}
}

func showMapConcurrentOperation() {
	m := map[int]int{
		1: 11,
		2: 12,
		3: 13,
	}

	go func() {
		for i := 0; i < 1000; i++ {
			doIteration(m)
		}
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			doWrite(m)
		}
	}()

	time.Sleep(5 * time.Second)
}

func doRead2(m map[int]string) {
	for k, v := range m {
		_ = fmt.Sprintf("k: %d, v: %s\n", k, v)
	}
}

func doWrite2(m map[int]string) {
	for k, v := range m {
		m[k] = fmt.Sprintf("v%s", v)
	}
}

func showMapConcurrentOperation2() {
	m := map[int]string{
		1: "1",
		2: "2",
		3: "3",
	}
	go func() {
		for i := 0; i < 1000; i++ {
			doRead2(m)
		}
	}()

	go func() {
		for i := 0; i < 1000; i++ {
			doWrite2(m)
		}
	}()
	// for i := 0; i < 100; i++ {
	// 	go doWrite2(m)
	// }

	time.Sleep(5)
}

func main() {
	// showMapConcurrentOperation()
	showMapConcurrentOperation2()
}
