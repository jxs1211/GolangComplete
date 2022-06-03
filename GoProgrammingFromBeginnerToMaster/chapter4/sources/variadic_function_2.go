package main

import "fmt"

func dump(args ...interface{}) {
	for _, v := range args {
		fmt.Println(v)
	}
}

func showVariadicTypeMatch() {
	//s := []string{"Tony", "John", "Jim"}
	s := []interface{}{"Tony", "John", "Jim"}
	dump(s...)
}

func Dump(args ...interface{}) {
	for i, v := range args {
		fmt.Println(i, v)
	}
}

type S struct {
	f int
}

func main() {
	Dump([]interface{}{"shen", "xian", "jie"}...)
	Dump([]interface{}{1, 2, 3}...)
	Dump([]interface{}{S{f: 1}, S{f: 2}, S{f: 3}}...)
}
