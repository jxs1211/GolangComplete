package main

import (
	"fmt"
	"sync"
	"time"
)

func newNumGenerator(start, count int) <-chan int {
	c := make(chan int)
	go func() {
		for i := start; i < start+count; i++ {
			c <- i
		}
		close(c)
	}()
	return c
}

func filterOdd(in int) (int, bool) {
	if in%2 != 0 {
		return 0, false
	}
	return in, true
}

func square(in int) (int, bool) {
	return in * in, true
}

func spawnGroup(name string, num int, f func(int) (int, bool), in <-chan int) <-chan int {
	groupOut := make(chan int)
	var outSlice []chan int
	for i := 0; i < num; i++ {
		out := make(chan int)
		go func(i int) {
			name := fmt.Sprintf("%s-%d:", name, i)
			fmt.Printf("%s begin to work...\n", name)

			for v := range in {
				r, ok := f(v)
				if ok {
					out <- r
				}
			}
			close(out)
			fmt.Printf("%s work done\n", name)
		}(i)
		outSlice = append(outSlice, out)
	}

	// Fan-in
	//
	// out --\
	//        \
	// out ---- --> groupOut
	//        /
	// out --/
	//
	go func() {
		var wg sync.WaitGroup
		for _, out := range outSlice {
			wg.Add(1)
			go func(out <-chan int) {
				for v := range out {
					groupOut <- v
				}
				wg.Done()
			}(out)
		}
		wg.Wait()
		close(groupOut)
	}()

	return groupOut
}

func showConcurrency9() {
	in := newNumGenerator(1, 20)
	out := spawnGroup("square", 2, square, spawnGroup("filterOdd", 3, filterOdd, in))

	time.Sleep(3 * time.Second)

	for v := range out {
		fmt.Println(v)
	}
}

func ChannelNumsGenerator(start, count int) <-chan int {
	out := make(chan int)
	go func() {
		for i := start; i <= count; i++ {
			out <- i
		}
		close(out)
	}()
	return out
}

func chan2slice(chan int) []int {

}

func spawn91(name string, nums int, f func(i int) (int, bool), in <-chan int) <-chan int {
	gOut := make(chan int)
	// outSlice := make([]chan int, 0, nums)
	var outSlice []chan int
	var wg sync.WaitGroup
	// fan out
	// go func() {
	for i := 0; i < nums; i++ {
		out := make(chan int)
		go func(id int) {
			defer fmt.Printf("%s: fan out goroutine-%d done\n", name, id)
			for v := range in {
				v1, ok := f(v)
				if ok {
					out <- v1
				}
			}
			close(out)
		}(i)
		outSlice = append(outSlice, out)
		// outSlice[i] = out
	}
	// }()
	// fan in
	go func() {
		for i, out := range outSlice {
			wg.Add(1)
			go func(id int) {
				defer func() {
					wg.Done()
					fmt.Printf("%s: fan in goroutinue-%d done\n", name, id)
				}()
				for v := range out {
					gOut <- v
				}
			}(i)
		}
		wg.Wait()
		close(gOut)
	}()

	return gOut
}

func showConcurrency91() {
	in := ChannelNumsGenerator(1, 20)
	out := spawn91("square", 2, square, spawn91("filterOdd", 3, filterOdd, in))

	for v := range out {
		fmt.Println(v)
	}

}

func main() {
	// showConcurrency9()
	showConcurrency91()
}
