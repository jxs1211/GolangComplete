package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func showReader() {
	r := strings.NewReader("hello, gopher!\n")
	lr := io.LimitReader(r, 4)
	if _, err := io.Copy(os.Stdout, lr); err != nil {
		log.Fatal(err)
	}
}

func showReader2() {
	r := strings.NewReader("Hello, world")
	lr := io.LimitReader(r, 10)
	if _, err := io.Copy(os.Stdout, lr); err != nil {
		fmt.Println(err)
	}
}

func main() {
	showReader2()
}
