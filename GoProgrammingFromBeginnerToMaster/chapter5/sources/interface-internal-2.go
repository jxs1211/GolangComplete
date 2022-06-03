package main

import "fmt"

type T int

func (t T) Error() string {
	return "bad error"
}

func printNilInterface() {
	// nil接口变量
	var i interface{} // 空接口类型
	var err error     // 非空接口类型
	println(i)
	println(err)
	println("i = nil:", i == nil)
	println("err = nil:", err == nil)
	println("i = err:", i == err)
	println("")
}

func printEmptyInterface() {
	// empty接口变量
	var eif1 interface{} // 空接口类型
	var eif2 interface{} // 空接口类型
	var n, m int = 17, 18

	eif1 = n
	eif2 = m

	println("eif1:", eif1)
	println("eif2:", eif2)
	println("eif1 = eif2:", eif1 == eif2)

	eif2 = 17
	println("eif1:", eif1)
	println("eif2:", eif2)
	println("eif1 = eif2:", eif1 == eif2)

	eif2 = int64(17)
	println("eif1:", eif1)
	println("eif2:", eif2)
	println("eif1 = eif2:", eif1 == eif2)

	println("")
}

func printNonEmptyInterface() {
	var err1 error // 非空接口类型
	var err2 error // 非空接口类型
	err1 = (*T)(nil)
	println("err1:", err1)
	println("err1 = nil:", err1 == nil)

	err1 = T(5)
	err2 = T(6)
	println("err1:", err1)
	println("err2:", err2)
	println("err1 = err2:", err1 == err2)

	err2 = fmt.Errorf("%d\n", 5)
	println("err1:", err1)
	println("err2:", err2)
	println("err1 = err2:", err1 == err2)

	println("")
}

func printEmptyInterfaceAndNonEmptyInterface() {
	var eif interface{} = T(5)
	var err error = T(5)
	println("eif:", eif)
	println("err:", err)
	println("eif = err:", eif == err)

	err = T(6)
	println("eif:", eif)
	println("err:", err)
	println("eif = err:", eif == err)
}

func printNilInterface2() {
	var efa interface{}
	var ifa error
	println(efa)
	println(ifa)

	println("efa == nil:", efa == nil)
	println("ifa == nil:", ifa == nil)
	println("")
}

func printEmptyInterface2() {
	var efa1 interface{}
	var efa2 interface{}
	m, n := 7, 8
	efa1, efa2 = m, n
	println("efa1:", efa1)
	println("efa2:", efa2)
	println("efa1 == efa2:", efa1 == efa2)

	efa2 = 7
	println("efa1:", efa1)
	println("efa2:", efa2)
	println("efa1 == efa2:", efa1 == efa2)

	efa2 = int64(7)
	println("efa1:", efa1)
	println("efa2:", efa2)
	println("efa1 == efa2:", efa1 == efa2)

	println("")
}

type T2 int

func (T2) Error() string {
	return fmt.Sprintf("%d")
}

func printNonEmptyInterface2() {
	var ifa1 error
	var ifa2 error
	ifa1 = (*T2)(nil)
	println(ifa1)
	println("ifa1 == nil", ifa1 == nil)

	ifa1, ifa2 = T2(5), T2(6)
	println("ifa1:", ifa1)
	println("ifa2:", ifa2)
	println("ifa1 == ifa2:", ifa1 == ifa2)

	ifa2 = T2(5)
	println("ifa1:", ifa1)
	println("ifa2:", ifa2)
	println("ifa1 == ifa2:", ifa1 == ifa2)

	println("")
}

func printEmptyInterfaceAndNonEmptyInterface2() {
	var efa interface{} = T2(5)
	var ifa error = T2(5)
	println("efa:", efa)
	println("ifa:", ifa)
	println("efa == ifa", efa == ifa)

	ifa = T2(6)
	println("efa:", efa)
	println("ifa:", ifa)
	println("efa == ifa", efa == ifa)

	println("")
}

func showPrintlnEfaceAndIfaceExample() {
	printNilInterface()
	printEmptyInterface()
	printNonEmptyInterface()
	printEmptyInterfaceAndNonEmptyInterface()
}

func showPrintlnEfaceAndIfaceExample2() {
	printNilInterface2()
	printEmptyInterface2()
	printNonEmptyInterface2()
	printEmptyInterfaceAndNonEmptyInterface2()
}

func main() {
	showPrintlnEfaceAndIfaceExample2()
}
