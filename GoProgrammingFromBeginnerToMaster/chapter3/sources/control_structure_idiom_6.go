package main

import (
	"fmt"
	"time"
)

func showBreakStopWhere() {
	exit := make(chan interface{})

	go func() {
		for {
			select {
			case <-time.After(time.Second):
				fmt.Println("tick")
			case <-exit:
				fmt.Println("exiting...")
				break
			}
		}
		fmt.Println("exit!")
	}()

	time.Sleep(3 * time.Second)
	exit <- struct{}{}

	// wait child goroutine exit
	time.Sleep(3 * time.Second)
}

func showBreakStopWhere2() {
	exit := make(chan int)

	go func() {
	outloop:
		for {
			select {
			case <-time.After(time.Second):
				fmt.Println("tick2")
				// outloop:
				for {
					select {
					case <-time.After(time.Second):
						fmt.Println("tick")
					case <-exit:
						fmt.Println("exiting")
						break outloop // jump to outloop position to continue
					}
				}
			}
		}
	}()
	time.Sleep(3 * time.Second)
	exit <- 1
	time.Sleep(3 * time.Second)
}

func showBreakAndContinueLabelWhere2() {
	// find the first element which is larger than 5 in every group
	sl := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}
	r := make([]int, 0, 9)
	for _, v := range sl {
	outloop:
		for _, e := range v {
			if e > 5 {
				r = append(r, e)
				break outloop
			}
		}
	}
	fmt.Println(r)

	for _, v := range sl {
	label:
		for _, e := range v {
			if e > 5 {
				r = append(r, e)
				continue label
			}
		}
	}
	fmt.Println(r)
}

func main() {
	// showBreakStopWhere2()
	showBreakAndContinueLabelWhere2()
}
