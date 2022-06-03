// inspired by https://github.com/go-functional/core
package main

import (
	"fmt"
)

type IntSliceFunctor interface {
	Fmap(fn func(int) int) IntSliceFunctor
}

type intSliceFunctorImpl struct {
	ints []int
}

func (isf intSliceFunctorImpl) Fmap(fn func(int) int) IntSliceFunctor {
	newInts := make([]int, len(isf.ints))
	for i, elt := range isf.ints {
		retInt := fn(elt)
		newInts[i] = retInt
	}
	return intSliceFunctorImpl{ints: newInts}
}

func NewIntSliceFunctor(slice []int) IntSliceFunctor {
	return intSliceFunctorImpl{ints: slice}
}

func showFunctor() {
	// 原切片
	intSlice := []int{1, 2, 3, 4}
	fmt.Printf("init a functor from int slice: %#v\n", intSlice)
	f := NewIntSliceFunctor(intSlice)
	fmt.Printf("original functor: %+v\n", f)

	mapperFunc1 := func(i int) int {
		return i + 10
	}

	mapped1 := f.Fmap(mapperFunc1)
	fmt.Printf("mapped functor1: %+v\n", mapped1)

	mapperFunc2 := func(i int) int {
		return i * 3
	}
	mapped2 := mapped1.Fmap(mapperFunc2)
	fmt.Printf("mapped functor2: %+v\n", mapped2)
	fmt.Printf("original functor: %+v\n", f) // 原functor没有改变
	fmt.Printf("composite functor: %+v\n", f.Fmap(mapperFunc1).Fmap(mapperFunc2))
}

type sliceFunctor interface {
	Fmap2(func(int, int) int) sliceFunctor
}

type sliceFunctorContainer struct {
	// c []int
	c <-chan int
}

func (sfc sliceFunctorContainer) Fmap2(f func(int, int) int) sliceFunctor {
	ch := make(chan int)
	go func() {
		for v := range sfc.c {
			ch <- f(v, v)
		}
		close(ch)
	}()
	return sliceFunctorContainer{c: ch}
}

func do(f func(int, int) int, in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for v := range in {
			out <- f(v, v)
		}
		close(out)
	}()
	return out
}

func showFunctor2() {
	// sig := make(chan struct{})
	rawSlice := []int{1, 2, 3, 4}
	in := make(chan int)
	go func() {
		for _, v := range rawSlice {
			in <- v
		}
		close(in)
	}()
	add := func(x, y int) int { return x + y }
	multiply := func(x, y int) int { return x * y }

	sfc := sliceFunctorContainer{c: in}
	mapped := sfc.Fmap2(add).Fmap2(multiply)
	sfc2 := (mapped).(sliceFunctorContainer)
	fmt.Printf("mapped functor2: %+v\n", mapped)
	fmt.Printf("original functor: %+v\n", sfc) // 原functor没有改变
	fmt.Printf("composite functor: %+v, sl: %v\n", sfc2, chan2slice(sfc2.c))
}

func chan2slice(ch <-chan int) []int {
	sl := make([]int, len(ch))
	for v := range ch {
		sl = append(sl, v)
	}
	return sl
}

func main() {
	showFunctor2()
	showFunctor()
}
