package main

import (
	"fmt"

	_ "github.com/jxs1211/package-init-order2/pkg1"
	"github.com/jxs1211/package-init-order2/pkg3"
	// "github.com/jxs1211/package-init-order2/pkg4"
)

var (
	_  = checkInit()
	v1 = initVariable("v1")
	v2 = initVariable("v2")
)

const (
	c1 = "c1"
	c2 = "c2"
)

func initVariable(s string) string {
	fmt.Printf("main init %s\n", s)
	return s
}

func checkInit() string {
	if c1 != "" {
		fmt.Println("main init c1")
	}
	if c2 != "" {
		fmt.Println("main init c2")
	}
	return ""
}

func init() {
	fmt.Println("main init")
}

func main() {
	// pkg1.ExportFunc()
	fmt.Println("main main")
	fmt.Println("pkg3 variable: ", pkg3.V1)
	pkg3.ExportFunc()
	// pkg3.inner()
	// pkg4.Print("shen")
}
