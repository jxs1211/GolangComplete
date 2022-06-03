package main

import (
	"fmt"
	"runtime"
	"time"
)

func deadloop() {
	for {
	}
}

func main() {
	fmt.Println(runtime.NumCPU())
	go deadloop()
	for {
		time.Sleep(time.Second * 1)
		fmt.Println("I got scheduled!")
	}
}
