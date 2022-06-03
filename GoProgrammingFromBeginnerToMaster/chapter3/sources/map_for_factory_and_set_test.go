package main

import (
	"fmt"
	"testing"
)

type Set map[int]bool

type Factory struct {
	m map[string]func(op int) int
}

func (f Factory) GetFunc(fn string) (func(int) int, error) {
	v, ok := f.m[fn]
	if !ok {
		return nil, fmt.Errorf("no func found: %s\n", fn)
	}
	return v, nil
}

func TestMapForFuncFactory(t *testing.T) {
	f := Factory{
		m: map[string]func(op int) int{
			"oneTimes":   func(op int) int { return op },
			"twoTimes":   func(op int) int { return op * 2 },
			"threeTimes": func(op int) int { return op * 3 },
		},
	}
	oneTimes, twoTimes, threeTimes := "oneTimes", "twoTimes", "threeTimes"
	fn, err := f.GetFunc(oneTimes)
	if err != nil {
		t.Log(err)
	}
	fmt.Println(fn(2))
	fn, err = f.GetFunc(twoTimes)
	if err != nil {
		t.Log(err)
	}
	fmt.Println(fn(2))
	fn, err = f.GetFunc(threeTimes)
	if err != nil {
		t.Log(err)
	}
	fmt.Println(fn(2))
}

func TestMapForset(t *testing.T) {
	// use map[type]bool as set
	// set := map[int]bool{}
	set := Set{}
	set[1] = true // create or update
	n := 1
	if set[n] {
		t.Logf("%d exist in set: %v\n", n, set)
	}
	n = 3
	if set[n] { // get
		t.Logf("%d exist in set: %v\n", n, set)
	} else {
		t.Logf("%d doesn't exist in set: %v\n", n, set)
	}
	n = 1
	delete(set, n) // delete
	if set[n] {
		t.Logf("%d exist in set: %v\n", n, set)
	} else {
		t.Logf("%d doesn't exist in set: %v\n", n, set)
	}
}
