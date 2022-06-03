package array_test

import "testing"

var a [100]int
var sl = a[:]

func arrayRangeLoop() {
	var sum int
	for _, n := range a {
		sum += n
	}
}

func pointerToArrayRangeLoop() {
	var sum int
	for _, n := range &a {
		sum += n
	}
}

func sliceRangeLoop() {
	var sum int
	for _, n := range sl {
		sum += n
	}
}

func arrayRangeLoop2() {
	var a [100]int
	for i, v := range a {
		_, _ = i, v
	}
}

func pointerToArrayRangeLoop2() {
	var a [100]int
	for i, v := range &a {
		_, _ = i, v
	}
}

func sliceRangeLoop2() {
	s := make([]int, 100)
	for i, v := range s {
		_, _ = i, v
	}
}

func BenchmarkArrayRangeLoop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		arrayRangeLoop()
	}
}
func BenchmarkPointerToArrayRangeLoop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pointerToArrayRangeLoop()
	}
}
func BenchmarkSliceRangeLoop(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sliceRangeLoop()
	}
}
