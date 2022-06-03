package main

import (
	"fmt"
	"time"
)

func producer(c chan<- int, sign chan struct{}) chan struct{} {
	var i int = 1
	s := make(chan struct{})
	go func(i int) {
		for {
			select {
			case <-sign:
				fmt.Println("receiver already exited")
				s <- struct{}{}
				return
			default:
				time.Sleep(2 * time.Second)
				ok := trySend(c, i)
				if !ok {
					fmt.Printf("[producer]: try send [%d], but channel is full\n", i)
					// sign <- struct{}{}
					close(sign)
					return
				}
				fmt.Printf("[producer]: send [%d] to channel\n", i)
				i++
				continue
			}

		}
	}(i)
	return s
}

func tryRecv(c <-chan int) (int, bool) {
	select {
	case i := <-c:
		return i, true

	default:
		return 0, false
	}
}

func trySend(c chan<- int, i int) bool {
	select {
	case c <- i:
		return true
	default:
		return false
	}
}

func consumer(c <-chan int, sign chan struct{}) {
	for {
		i, ok := tryRecv(c)
		if !ok {
			fmt.Println("[consumer]: try to recv from channel, but the channel is empty")
			time.Sleep(1 * time.Second)
			continue
		}
		fmt.Printf("[consumer]: recv [%d] from channel\n", i)
		if i >= 3 {
			fmt.Println("[consumer]: exit")
			close(sign)
			return
		}
	}
}

func main() {
	c := make(chan int, 3)
	sign := make(chan struct{})
	s := producer(c, sign)
	go consumer(c, sign)

	<-s
	println("done")
	// select {} // 故意阻塞在此
}
