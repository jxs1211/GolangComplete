package main

import (
	"fmt"
	"net/http"
	"reflect"
	"sync"
)

type Node struct {
	left  *Node
	right *Node
}

type Person struct {
	name string
	age  int
}

func main() {
	http.HandleFunc
	// p := &Person{}
	var p Person
	// if p != nil {
	// 	fmt.Printf("The Person: %+v\n", p)
	// 	return
	// }
	// if p == struct{}{} {
	// 	fmt.Println("equal")
	// }
	p2 := Person{}
	if ok := reflect.DeepEqual(p, p2); ok {
		fmt.Println("equal")
	} else {
		fmt.Println("not equal")
	}
	fmt.Printf("Person: %+v\n", p)
	fmt.Printf("name: %q\n", p.name)
}

func showGoroutinue() {
	var wg sync.WaitGroup
	done := make(chan struct{})
	wq := make(chan interface{})
	workerCount := 2

	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go doit(i, wq, done, &wg)
	}

	for i := 0; i < workerCount; i++ {
		wq <- i
	}

	close(done)
	wg.Wait()
	fmt.Println("all done!")
}

func doit(workerId int, wq <-chan interface{}, done <-chan struct{}, wg *sync.WaitGroup) {
	fmt.Printf("[%v] is running\n", workerId)
	defer wg.Done()
	for {
		select {
		case m := <-wq:
			fmt.Printf("[%v] m => %v\n", workerId, m)
		case <-done:
			fmt.Printf("[%v] is done\n", workerId)
			return
		}
	}
}

func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

func twoSum(nums []int, target int) []int {
	m := make(map[int]int, len(nums))
	for i, v := range nums {
		_diff := target - v
		if index, ok := m[_diff]; ok {
			return []int{index, i}
		}
		m[v] = i
	}
	return []int{}
}

func searchInsert(nums []int, target int) int {
	l, r := 0, len(nums)-1
	for l <= r {
		mid := l + (r-l)>>1
		if nums[mid] >= target {
			r = mid - 1
		} else {
			if mid == len(nums)-1 || nums[mid+1] >= target {
				return mid + 1
			}
			l = mid + 1
		}
	}
	return 0
}

func Add(s []int, i, v int) []int {
	return s
}

func Delete(s []int, i, v int) []int {
	return s
}
