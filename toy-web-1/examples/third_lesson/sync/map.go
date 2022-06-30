package main

import (
	"fmt"
	"sync"
)

func showSyncMap() {
	m := sync.Map{}
	m.Store("cat", "Tom")
	m.Store("dog", "Donald")

	val, ok := m.Load("cat")
	if ok {
		fmt.Println(len(val.(string)))
	}

}

func main() {
	showSyncMap()
}

func show() {
	m := sync.Map{}
	m.Store("cat", "Tom")
	m.Store("mouse", "Jerry")

	// 这里重新读取出来的，就是
	// val是接口，需要自己转换
	val, ok := m.Load("mouse")
	if ok {
		fmt.Println(len(val.(string)))
	}

}
