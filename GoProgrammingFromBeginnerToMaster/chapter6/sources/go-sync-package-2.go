package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type foo struct {
	n int
	sync.Mutex
}

func CopySyncType() {
	f := foo{n: 17}

	go func(f foo) {
		for {
			log.Println("g2: try to lock foo...")
			f.Lock()
			log.Println("g2: lock foo ok")
			time.Sleep(3 * time.Second)
			f.Unlock()
			log.Println("g2: unlock foo ok")
		}
	}(f)

	f.Lock()
	log.Println("g1: lock foo ok")

	// 在mutex首次使用后复制其值
	go func(f foo) {
		for {
			log.Println("g3: try to lock foo...")
			f.Lock()
			log.Println("g3: lock foo ok")
			time.Sleep(5 * time.Second)
			f.Unlock()
			log.Println("g3: unlock foo ok")
		}
	}(f)

	time.Sleep(1000 * time.Second)
	f.Unlock()
	log.Println("g1: unlock foo ok")
}

type clock struct {
	i int
	sync.Mutex
}

func CopySyncType2() {
	s := clock{i: 1}

	go func(c clock) {
		for {
			fmt.Println("g1 try lock...")
			c.Lock()
			fmt.Println("g1 lock")
			time.Sleep(3 * time.Second)
			c.Unlock()
			fmt.Println("g1 unlock")
		}
	}(s)

	s.Lock()
	fmt.Println("g0 lock")
	go func(c clock) {
		for {
			fmt.Println("g2 try lock...")
			c.Lock()
			fmt.Println("g2 lock")
			time.Sleep(5 * time.Second)
			c.Unlock()
			fmt.Println("g2 unlock")
		}
	}(s)

	time.Sleep(1000 * time.Second)
	s.Unlock()
	fmt.Println("g0 unlock")
}

func main() {
	// CopySyncType()
	CopySyncType2()
}
