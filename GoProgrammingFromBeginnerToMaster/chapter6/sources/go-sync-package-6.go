package main

import (
	"log"
	"sync"
	"time"
)

type Foo struct {
}

var once sync.Once
var instance *Foo

func GetInstance(id int) *Foo {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("goroutine-%d: caught a panic: %s", id, e)
		}
	}()
	log.Printf("goroutine-%d: enter GetInstance\n", id)
	once.Do(func() {
		instance = &Foo{}
		time.Sleep(3 * time.Second)
		log.Printf("goroutine-%d: the addr of instance is %p\n", id, instance)
		panic("panic in once.Do function")
	})
	return instance
}

func showOnceDo() {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(i int) {
			inst := GetInstance(i)
			log.Printf("goroutine-%d: the addr of instance returned is %p\n", i, inst)
			wg.Done()
		}(i + 1)
	}
	time.Sleep(5 * time.Second)
	inst := GetInstance(0)
	log.Printf("goroutine-0: the addr of instance returned is %p\n", inst)

	wg.Wait()
	log.Printf("all goroutines exit\n")
}

type inst struct {
}

var ist *inst

// var instance *Foo
var once2 sync.Once

func getInst(id int) *inst {
	log.Printf("goroutine-%d: enter getInst", id)
	defer func() {
		if p := recover(); p != nil {
			log.Println("caught a panic: ", p)
		}
	}()
	once2.Do(func() {
		ist = &inst{}
		log.Printf("goroutine-%d: getInst: %p\n", id, ist)
		panic("raise a panic")
	})
	return ist
}

func showOnceDo2() {
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			ist := getInst(id)
			log.Printf("goroutine-%d: instance: %+v %p\n", id, ist, ist)
		}(i + 1)
	}
	time.Sleep(5 * time.Second)
	ist := getInst(0)
	log.Printf("goroutine-%d: get instance: %p\n", 0, ist)

	wg.Wait()
	log.Println("done")
}

func main() {
	// showOnceDo()
	showOnceDo2()
}
