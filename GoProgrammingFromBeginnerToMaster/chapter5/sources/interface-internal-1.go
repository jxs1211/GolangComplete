package main

import (
	"errors"
	"fmt"
)

type MyError struct {
	error
}

var ErrBad = MyError{
	error: errors.New("bad error"),
}

func bad() bool {
	return false
}

func returnsError() error {
	var p *MyError = nil
	if bad() {
		p = &ErrBad
	}
	return p
}

func showInterfaceAndNilCompare() {
	e := returnsError()
	if e != nil {
		fmt.Printf("error: %+v\n", e)
		return
	}
	fmt.Println("ok")
}

type MyError2 struct {
	error
}

var ErrBad2 = MyError2{
	error: errors.New("bad error"),
}

func bad2() bool { return false }

func returnError2() error {
	var e *MyError2 = nil
	if bad2() {
		e = &ErrBad2
	}
	return e
}

func showInterfaceAndNilCompare2() {
	err := returnError2()
	if v, ok := err.(*MyError2); ok {
		if v != nil {
			// fmt.Println(err)
			println(err)
			return
		}
	}
	fmt.Println("ok")
}

func main() {
	showInterfaceAndNilCompare2()
}
