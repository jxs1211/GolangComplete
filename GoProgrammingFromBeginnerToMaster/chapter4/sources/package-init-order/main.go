package main

import (
	"fmt"
	"log"

	_ "github.com/bigwhite/package-init-order/pkg1"
	_ "github.com/bigwhite/package-init-order/pkg3"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/gorm"
)

var (
	_  = constInitCheck()
	v1 = variableInit("v1")
	v2 = variableInit("v2")
)

const (
	c1 = "c1"
	c2 = "c2"
)

func constInitCheck() string {
	if c1 != "" {
		fmt.Println("main: const c1 init")
	}
	if c2 != "" {
		fmt.Println("main: const c2 init")
	}
	return ""
}

func variableInit(name string) string {
	fmt.Printf("main: var %s init\n", name)
	return name
}

func init() {
	fmt.Println("main: init")
}

func main() {
	// do nothing
	db, err := gorm.Open("mysql", "root:stadmin@(localhost)/bitstormx?charset=utf8&parseTime=True&loc=Local")

	if err != nil {
		log.Fatal(err)
	}
	rows, err := db.Query("select * from test limit 1")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(rows)

}
