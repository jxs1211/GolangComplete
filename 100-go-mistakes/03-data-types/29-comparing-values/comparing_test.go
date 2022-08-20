package main

import (
	"fmt"
	"reflect"
	"testing"
)

func TestComparing(t *testing.T) {

	cust11 := customer1{id: "x"}
	cust12 := customer1{id: "x"}
	fmt.Println(cust11 == cust12)

	cust21 := customer2{id: "x", operations: []float64{1.}}
	cust22 := customer2{id: "x", operations: []float64{1.}}
	// Doesn't compile
	// fmt.Println(cust21 == cust22)
	_ = cust21
	_ = cust22

	var a any = 3
	var b any = 3
	fmt.Println(a == b)

	var cust31 any = customer2{id: "x", operations: []float64{1.}}
	var cust32 any = customer2{id: "x", operations: []float64{1.}}
	fmt.Println(cust31 == cust32)

	cust41 := customer2{id: "x", operations: []float64{1.}}
	cust42 := customer2{id: "x", operations: []float64{1.}}
	fmt.Println(reflect.DeepEqual(cust41, cust42))

}
