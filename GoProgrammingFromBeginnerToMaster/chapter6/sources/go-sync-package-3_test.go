package main

import (
	"sync"
	"testing"
)

var data1 int
var mut1 sync.Mutex
var data2 int
var mut2 sync.RWMutex

func BenchmarkSyncReadByMutex(b *testing.B) {
	b.RunParallel(func(b *testing.PB) {
		for b.Next() {
			mut1.Lock()
			_ = data1
			mut1.Unlock()
		}
	})
}

func BenchmarkSyncReadByRWMutexRLock(b *testing.B) {
	b.RunParallel(func(b *testing.PB) {
		for b.Next() {
			mut2.RLock()
			_ = data1
			mut2.RUnlock()
		}
	})
}

func BenchmarkSyncReadByRWMutexLock(b *testing.B) {
	b.RunParallel(func(b *testing.PB) {
		for b.Next() {
			mut2.Lock()
			_ = data1
			mut2.Unlock()
		}
	})
}

func BenchmarkSyncWriteByMutex(b *testing.B) {
	b.RunParallel(func(b *testing.PB) {
		for b.Next() {
			mut1.Lock()
			data1++
			mut1.Unlock()
		}
	})
}

func BenchmarkSyncWriteByRWMutexRLock(b *testing.B) {
	b.RunParallel(func(b *testing.PB) {
		for b.Next() {
			mut2.RLock()
			data1++
			mut2.RUnlock()
		}
	})
}

func BenchmarkSyncWriteByRWMutexLock(b *testing.B) {
	b.RunParallel(func(b *testing.PB) {
		for b.Next() {
			mut2.Lock()
			data1++
			mut2.Unlock()
		}
	})
}
