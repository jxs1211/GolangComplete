package main

import "testing"

const length = 100000

type S5 struct {
	i1, i2, i3, i4, i5 int
}

type S4 struct {
	i1, i2, i3, i4 int
}

type S3 struct {
	i1, i2, i3 int
}

type S2 struct {
	i1, i2 int
}

type S1 struct {
	i1 int
}

func BenchmarkForLoopLargeArrayWithStructU5Element(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var a [length]S5
		for j := 0; j < length; j++ {
			a[j] = S5{i1: j}
		}
	}
}

func BenchmarkForLoopLargeArrayWithStructU4Element(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var a [length]S4
		for j := 0; j < length; j++ {
			a[j] = S4{i1: j}
		}
	}
}

func BenchmarkForLoopLargeArrayWithStructU3Element(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var a [length]S3
		for j := 0; j < length; j++ {
			a[j] = S3{i1: j}
		}
	}
}

func BenchmarkForLoopLargeArrayWithStructU2Element(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var a [length]S2
		for j := 0; j < length; j++ {
			a[j] = S2{i1: j}
		}
	}
}

func BenchmarkForLoopLargeArrayWithStructU1Element(b *testing.B) {
	b.ReportAllocs()
	var a [length]S1
	for i := 0; i < b.N; i++ {
		for j := 0; j < length; j++ {
			a[j] = S1{i1: j}
		}
	}
}

func BenchmarkASForLoopLargeArrayWithStructU5Element(b *testing.B) {
	b.ReportAllocs()
	var a [length]S5
	for i := 0; i < b.N; i++ {
		for j := 0; j < length; j++ {
			a[j].i1 = j
		}
	}
}

func BenchmarkASForLoopLargeSliceWithStructU5Element(b *testing.B) {
	b.ReportAllocs()
	a := make([]S5, length)
	for i := 0; i < b.N; i++ {
		for j := 0; j < length; j++ {
			// a[j].i1 = j
			a[j].i1 = j
		}
	}
}

func BenchmarkASForLoopLargeSliceWithStructU5Element2(b *testing.B) {
	b.ReportAllocs()
	a := make([]S5, length)
	for i := 0; i < b.N; i++ {
		for j := 0; j < length; j++ {
			// for j, _ := range a {
			a[j].i1 = j
			// a = append(a, a[j])
		}
	}
}

func BenchmarkASForRangeLargeArrayWithStructU5Element(b *testing.B) {
	b.ReportAllocs()
	var a [length]S5
	for i := 0; i < b.N; i++ {
		for j, _ := range a {
			// a[j] = S5{i1: j}
			a[j].i1 = j
		}
	}
}

func BenchmarkASForRangeLargeSliceWithStructU5Element(b *testing.B) {
	b.ReportAllocs()
	a := make([]S5, length)
	for i := 0; i < b.N; i++ {
		for j, _ := range a {
			// a[j] = S5{i1: j}
			a[j].i1 = j
		}
	}
}

func BenchmarkASForRangeLargeSliceWithStructU5Element2(b *testing.B) {
	b.ReportAllocs()
	a := make([]S5, 0, length)
	for i := 0; i < b.N; i++ {
		for j, _ := range a {
			// a[j] = S5{i1: j}
			a[j].i1 = j
		}
	}
}

func TestForLoopAndForRange(t *testing.T) {
	a := make([]S5, 0, length)
	t.Log(a, len(a))
	for j, _ := range a {
		t.Log(j)
		// for j := 0; j < length; j++ {
		a[j].i1 = j
	}
}

func BenchmarkASForRangeLargeArrayWithStructU4Element(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var a [length]S4
		for j, _ := range a {
			a[j] = S4{i1: j}
		}
	}
}

func BenchmarkForRangeLargeArrayWithStructU3Element(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var a [length]S3
		for j, _ := range a {
			a[j] = S3{i1: j}
		}
	}
}

func BenchmarkForRangeLargeArrayWithStructU2Element(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var a [length]S2
		for j, _ := range a {
			a[j] = S2{i1: j}
		}
	}
}

func BenchmarkForRangeLargeArrayWithStructU1Element(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var a [length]S1
		for j, _ := range a {
			a[j] = S1{i1: j}
		}
	}
}

func BenchmarkForLoopLargeArrayWithIntElement(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var a [length]int
		for j := 0; j < length; j++ {
			a[j] = j
		}
	}
}

func BenchmarkForRangeLargeArrayWithIntElement(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var a [length]int
		for j, _ := range a {
			a[j] = j
		}
	}
}
