package defer_test

import "testing"

func sum(max int) int {
	total := 0
	for i := 0; i < max; i++ {
		total += i
	}

	return total
}

func fooWithDefer() {
	defer func() {
		sum(10)
	}()
}
func fooWithoutDefer() {
	sum(10)
}

// func BenchmarkFooWithDefer(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		fooWithDefer()
// 	}
// }
// func BenchmarkFooWithoutDefer(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		fooWithoutDefer()
// 	}
// }

func forLoop() {
	for i := 0; i < 100; i++ {
		_ = i
	}
}

func forLoopWithPanic() {
	defer func() {
		_ = recover()
	}()
	for i := 0; i < 100; i++ {
		_ = i
		panic(-1)
	}
}

func FuncWithPanicRecover() {
	forLoopWithPanic()
}

func FuncWithDefer() {
	defer func() {
		forLoop()
	}()
}

func FuncWithoutDefer() {
	forLoop()
}

func BenchmarkFuncWithDefer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FuncWithDefer()
	}
}

func BenchmarkFuncWithoutDefer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FuncWithoutDefer()
	}
}

// func BenchmarkFuncWithPanicRecover(b *testing.B) {
// 	for i := 0; i < b.N; i++ {
// 		FuncWithPanicRecover()
// 	}
// }
