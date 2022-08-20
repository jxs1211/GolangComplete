package main

import (
	"runtime"
	"testing"
)

func TestLeak(t *testing.T) {
	// Init
	n := 1_000_000
	m := make(map[int][128]byte)
	printAlloc()

	// Add elements
	for i := 0; i < n; i++ {
		m[i] = randBytes()
	}
	printAlloc()

	// Remove elements
	for i := 0; i < n; i++ {
		delete(m, i)
	}

	// End
	runtime.GC()
	printAlloc()
	runtime.KeepAlive(m)
}
