package main

import (
	"fmt"
	"sync"
	"time"
)

var mutex sync.Mutex
var rwMutex sync.RWMutex

func Mutex() {
	mutex.Lock()
	defer mutex.Unlock()
	// 你的代码
}

// 尽量使用RWMutex，不要使用mutex，
func RwMutex() {
	// 加读锁
	// 一定要通过defer解锁，否则panic就无法解锁了s
	rwMutex.RLock()
	defer rwMutex.RUnlock()

	// 也可以加写锁
	rwMutex.Lock()
	defer rwMutex.Unlock()
}

// 不可重入例子
func Failed1() {
	mutex.Lock()
	defer mutex.Unlock()

	// 这一句会死锁
	// 但是如果你只有一个goroutine，那么这一个会导致程序崩溃
	mutex.Lock()
	defer mutex.Unlock()
}

// 不可升级
func Failed2() {
	rwMutex.RLock()
	defer rwMutex.RUnlock()

	// 这一句会死锁
	// 但是如果你只有一个goroutine，那么这一个会导致程序崩溃
	mutex.Lock()
	defer mutex.Unlock()
}

func main() {
	// 在子goroutine中运行Failed1是不影响主goroutine的，不会抛出panic, 只有一个goroutinue的情况下，就会panic
	go Failed1()
	time.Sleep(2 * time.Second)
	fmt.Println("Done")
}
