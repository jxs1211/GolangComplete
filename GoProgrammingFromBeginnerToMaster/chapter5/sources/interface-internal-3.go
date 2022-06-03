package main

import "unsafe"

type T int

func (t T) Error() string {
	return "bad error"
}

func showEfaceAndIfaceTypeVariablesComparison() {
	var eif interface{} = T(5)
	var err error = T(5)
	println("eif:", eif)
	println("err:", err)
	println("eif = err:", eif == err)

	dumpEface(eif)
	// data: bad error
	// _type: {size:8 ptrdata:0 hash:1156555957 tflag:15 align:8 fieldalign:8 kind:2 alg:0x49f2d0 gcdata:0x4b3f8a str:3100 ptrToThis:37760}

	// dumpItabOfIface(unsafe.Pointer(&err))
	// // _type: {size:8 ptrdata:0 hash:1156555957 tflag:15 align:8 fieldalign:8 kind:2 alg:0x49f2d0 gcdata:0x4b3f8a str:3100 ptrToThis:37760}
	// dumpDataOfIface(err)
	// // data: bad error
}

type T2 int

func (T2) Error() string { return "T2 bad error" }

func showEfaceAndIfaceTypeVariablesComparison2() {
	// ptrSize := unsafe.Sizeof(uintptr(0))
	// fmt.Println(ptrSize)
	var eif interface{} = T2(5)
	var iif error = T2(5)
	println("eif:", eif)
	println("iif:", iif)
	println("eif = iif:", eif == iif)

	dumpEface(eif)
	dumpItabOfIface(unsafe.Pointer(&iif))
	dumpDataOfIface(iif)
}

func main() {
	showEfaceAndIfaceTypeVariablesComparison2()
}
