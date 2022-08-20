package main

import (
	_ "embed"
	"plugin"
)

//go:generate go run ./template/gen_embed_var.go

func main() {
	p, err := plugin.Open("plugin.so")
	if err != nil {
		panic(err)
	}
	v, err := p.Lookup("V")
	if err != nil {
		panic(err)
	}
	f, err := p.Lookup("F")
	if err != nil {
		panic(err)
	}
	*v.(*int) = 8
	f.(func())() // prints "Hello, number 7"
}
