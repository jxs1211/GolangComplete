package main

import "fmt"

// Send the sequence 2, 3, 4, ... to channel 'ch'.
func Generate(ch chan<- int) {
	for i := 2; ; i++ {
		ch <- i // Send 'i' to channel 'ch'.
		fmt.Printf("Generate: send to in channel: %d\n", i)
	}
}

// Copy the values from channel 'in' to channel 'out',
// removing those divisible by 'prime'.
func Filter(in <-chan int, out chan<- int, prime int) {
	for {
		i := <-in // Receive value from 'in'.
		fmt.Printf("Filter: receive from in channel: %d\n", i)
		if i%prime != 0 {
			out <- i // Send 'i' to 'out'.
			fmt.Printf("Filter: send to out channel: %d\n", i)
		}
	}
}

// The prime sieve: Daisy-chain Filter processes.
func sieve() {
	ch := make(chan int) // Create a new channel.
	go Generate(ch)      // Launch Generate goroutine.
	for i := 0; i < 10; i++ {
		prime := <-ch
		fmt.Printf("main: receive from in channel: %d\n", prime)
		ch1 := make(chan int)
		go Filter(ch, ch1, prime)
		ch = ch1
	}
}

func main() {
	sieve2()
}

func generate2(ch chan int) {
	for i := 2; i < 10; i++ {
		ch <- i
		fmt.Printf("Generate: send to in channel: %d\n", i)
	}
}

func filter2(in <-chan int, out chan<- int, prime int) {
	for {
		i := <-in
		fmt.Printf("Filter: receive from in channel: %d\n", i)
		if i%prime != 0 {
			out <- i
			fmt.Printf("Filter: send to out channel: %d\n", i)
		}
	}
}

func sieve2() {
	ch := make(chan int)
	go generate2(ch)
	for i := 0; i < 3; i++ {
		prime := <-ch
		fmt.Printf("main: receive from in channel: %d\n", prime)
		out := make(chan int)
		go filter2(ch, out, prime)
		ch = out
	}
}
