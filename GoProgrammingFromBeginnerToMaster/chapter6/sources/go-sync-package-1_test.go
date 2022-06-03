package main

import (
	"sync"
	"testing"
)

var cs = 0 // 模拟临界区要保护的数据
var mu sync.Mutex
var c = make(chan struct{}, 1)

func criticalSectionSyncByMutex() {
	mu.Lock()
	cs++
	mu.Unlock()
}

func criticalSectionSyncByChan() {
	c <- struct{}{}
	cs++
	<-c
}

func BenchmarkCriticalSectionSyncByMutex(b *testing.B) {
	for n := 0; n < b.N; n++ {
		criticalSectionSyncByMutex()
	}
}

func BenchmarkCriticalSectionSyncByChan(b *testing.B) {
	for n := 0; n < b.N; n++ {
		criticalSectionSyncByChan()
	}
}

var i int
var j int
var mut sync.Mutex
var cha = make(chan struct{}, 1)

func CriticalSectionSyncByMutex2() {
	mut.Lock()
	defer mut.Unlock()
	_ = i
}

func CriticalSectionSyncByChan2() {
	cha <- struct{}{}
	_ = j
	<-cha
}

func CriticalSectionSyncWriteByMutex2() {
	mut.Lock()
	defer mut.Unlock()
	i++
}

func CriticalSectionSyncWriteByChan2() {
	cha <- struct{}{}
	j++
	<-cha
}

func BenchmarkCriticalSectionSyncWriteByMutex2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CriticalSectionSyncWriteByMutex2()
	}
}

func BenchmarkCriticalSectionSyncWriteByChan2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CriticalSectionSyncWriteByChan2()
	}
}

func BenchmarkCriticalSectionSyncByMutex2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CriticalSectionSyncByMutex2()
	}
}

func BenchmarkCriticalSectionSyncByChan2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		CriticalSectionSyncByChan2()
	}
}
