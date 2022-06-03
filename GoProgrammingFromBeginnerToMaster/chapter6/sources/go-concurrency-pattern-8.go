package main

import "fmt"

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

func spawn(f func(int) (int, bool), in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		for v := range in {
			r, ok := f(v)
			if ok {
				out <- r
			}
		}
		close(out)
	}()

	return out
}

func showConcurrencyPattern8() {
	in := newNumGenerator(1, 20)
	out := spawn(square, spawn(filterOdd, in))

	for v := range out {
		println(v)
	}
}

func ChannelGenerator(start, count int) <-chan int {
	out := make(chan int)
	go func() {
		for i := start; i <= count; i++ {
			out <- i
		}
		close(out)
	}()
	return out
}

func filterOdd81(i int) (int, bool) {
	if i%2 == 0 {
		return i, true
	}
	return 0, false
}

func square81(i int) (int, bool) {
	return i * i, true
}

func filterOver300(i int) (int, bool) {
	if i > 300 {
		return 0, false
	}
	return i, true
}

func spawn81(f func(i int) (int, bool), in <-chan int) <-chan int {
	out := make(chan int)

	go func() {
		for v := range in {
			v1, ok := f(v)
			if ok {
				out <- v1
			}
		}
		close(out)
	}()

	return out
}

// pipline mode
func showConcurrencyPattern81() {

	in := ChannelGenerator(1, 20)
	out := spawn81(filterOver300, spawn81(filterOdd81, spawn81(square81, in)))

	for v := range out {
		fmt.Println(v)
	}
	fmt.Println("all done")
}

func main() {
	// showConcurrency8()
	showConcurrencyPattern81()
}
