package chapter15

import (
	"testing"
)

const nums = 1000

type S struct {
	array [10000]int
}

func WithValue(a [nums]S) int {
	return 0
}

func WithPointer(a *[nums]S) int {
	return 0
}

func WithValue2(a [10000]int) int {
	return 0
}

func WithPointer2(a *[10000]int) int {
	return 0
}

func Test1(t *testing.T) {
	var a [nums]S
	t.Log(a[len(a)-1])
	t.Log(WithValue(a))
	t.Log(WithPointer(&a))
}

func BenchmarkWithValue(b *testing.B) {
	var a [nums]S
	for i := 0; i < b.N; i++ {
		_ = WithValue(a)
	}
}

func BenchmarkWithPointer(b *testing.B) {
	var a [nums]S
	for i := 0; i < b.N; i++ {
		_ = WithPointer(&a)
	}
}

func BenchmarkWithValue2(b *testing.B) {
	var a [10000]int
	for i := 0; i < b.N; i++ {
		_ = WithValue2(a)
	}
}

func BenchmarkWithPointer2(b *testing.B) {
	var a [10000]int
	for i := 0; i < b.N; i++ {
		_ = WithPointer2(&a)
	}
}

func TestSubtest(t *testing.T) {
	t.Run("case1", func(t *testing.T) { t.Log("Run case1") })
	t.Run("case2", func(t *testing.T) { t.Log("Run case2") })
	t.Run("case3", func(t *testing.T) { t.Log("Run case3") })
}
