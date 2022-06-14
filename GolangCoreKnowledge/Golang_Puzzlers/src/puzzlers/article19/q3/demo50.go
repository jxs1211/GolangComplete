package main

import (
	"errors"
	"fmt"
	"net/http"
)

func ShowPanicInRecover() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("fatal error: %s\n", p)
			panic("raise another error")
		}
	}()
	panic("something wrong")
}

func ShowDeferExample() {
	fmt.Println("Enter function main.")

	defer func() {
		fmt.Println("Enter defer function.")

		// recover函数的正确用法。
		if p := recover(); p != nil {
			fmt.Printf("panic: %s\n", p)
		}

		fmt.Println("Exit defer function.")
	}()

	// recover函数的错误用法。
	fmt.Printf("no panic: %v\n", recover())

	// 引发panic。
	panic(errors.New("something wrong"))

	// recover函数的错误用法。
	p := recover()
	fmt.Printf("panic: %s\n", p)

	fmt.Println("Exit function main.")
}

//
func TranslatePanicToError() (err error) {
	defer func() {
		if p := recover(); p != nil {
			err = fmt.Errorf("fatal error: %s", p)
		}
	}()
	panic("Ops!")
	return
}

func main() {
	// fmt.Println(utf8.RuneCountInString("shen沈"))
	// fmt.Println(len("shen沈"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi: %s\n", r.URL)
	})
	http.ListenAndServe(":8080", nil)
	// ShowPanicInRecover()
	// ShowDeferExample()
	// if err := TranslatePanicToError(); err != nil {
	// 	fmt.Println(err)
	// }
}
