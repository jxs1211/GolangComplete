package main

import "fmt"

func foo(b ...byte) {
	fmt.Println(string(b))
}

func showExceptionalExample() {
	b := []byte{}
	b = append(b, "hello"...)
	fmt.Println(string(b))

	// foo("hello"...)
}

func bar(args ...byte) {
	for i, v := range args {
		fmt.Println(i, v)
	}
}

func showExceptionalExample2() {
	var b []byte
	b = append(b, "shenxianjie"...)
	fmt.Println(string(b), len(b))

	// bar("shen"...)
}

func main() {
	showExceptionalExample2()
}
