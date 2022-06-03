package main

import (
	"fmt"
	"unsafe"
)

const ptrSize = unsafe.Sizeof(uintptr(0))

type typeFlag struct {
	hash  func(unsafe.Pointer, uintptr) uintptr
	equal func(unsafe.Pointer, unsafe.Pointer) bool
}

type tflag uint8
type nameOff int32
type typeOff int32

type _type struct {
	size       uintptr
	ptrdata    uintptr // size of memory prefix holding all pointers
	hash       uint32
	tflag      tflag
	align      uint8
	fieldAlign uint8
	kind       uint8
	// function for comparing objects of this type
	// (ptr to object A, ptr to object B) -> ==?
	equal func(unsafe.Pointer, unsafe.Pointer) bool
	// gcdata stores the GC type data for the garbage collector.
	// If the KindGCProg bit is set in kind, gcdata is a GC program.
	// Otherwise it is a ptrmask bitmap. See mbitmap.go for details.
	gcdata    *byte
	str       nameOff
	ptrToThis typeOff
}

type imethod struct {
	name nameOff
	ityp typeOff
}

type name struct {
	bytes *byte
}

type interfacetype struct {
	typ     _type
	pkgpath name
	mhdr    []imethod
}

type itab struct {
	inter *interfacetype
	_type *_type
	hash  uint32 // copy of _type.hash. Used for type switches.
	_     [4]byte
	fun   [1]uintptr // variable sized. fun[0]==0 means _type does not implement inter.
}

type iface struct {
	tab  *itab
	data unsafe.Pointer
}

type eface struct {
	_type *_type
	data  unsafe.Pointer
}

func dumpEface(i interface{}) {
	ptrEface := (*eface)(unsafe.Pointer(&i))

	if ptrEface._type != nil {
		fmt.Printf("\t type: %+v\n", *(ptrEface._type))
	}

	if ptrEface.data != nil {
		switch i.(type) {
		case int:
			dumpInt(ptrEface.data)
		case T2:
			dumpT2(ptrEface.data)
		case float64:
			dumpFloat64(ptrEface.data)
		default:
			fmt.Println("unsupport type")
		}
	}
}

// // fit for go 1.13.x version

func dumpItabOfIface(ptrToIface unsafe.Pointer) {
	p := (*iface)(ptrToIface)
	fmt.Printf("iface: %+v\n", *p)

	if p.tab != nil {
		// dump itab
		fmt.Printf("\t itab: %+v\n", *(p.tab))
		// dump inter in itab
		fmt.Printf("\t\t inter: %+v\n", *(p.tab.inter))

		// dump _type in itab
		fmt.Printf("\t\t _type: %+v\n", *(p.tab._type))

		// dump fun in tab
		funPtr := unsafe.Pointer(&(p.tab.fun))
		fmt.Printf("\t\t fun: [")
		for i := 0; i < len((*(p.tab.inter)).mhdr); i++ {
			tp := (*uintptr)(unsafe.Pointer(uintptr(funPtr) + uintptr(i)*ptrSize))
			fmt.Printf("0x%x(%d),", *tp, *tp)
		}
		fmt.Printf("]\n")
	}
}

func dumpDataOfIface(i interface{}) {
	// this is a trick as the data part of eface and iface are same
	ptrToEface := (*iface)(unsafe.Pointer(&i))

	if ptrToEface.data != nil {
		// dump data
		switch i.(type) {
		case int:
			dumpInt(ptrToEface.data)
		case float64:
			dumpFloat64(ptrToEface.data)
		case T:
			dumpT2(ptrToEface.data)

		// other cases ... ...

		default:
			fmt.Printf("\t data: unsupported type\n")
		}
	}
	fmt.Printf("\n")
}

func dumpT(dataOfIface unsafe.Pointer) {
	var p *T = (*T)(dataOfIface)
	fmt.Printf("\t data: %+v\n", *p)
}

func dumpT2(p unsafe.Pointer) {
	fmt.Printf("\t data: %+v\n", *(*T2)(p))
}

func dumpInt(p unsafe.Pointer) {
	fmt.Printf("\t data: %v\n", *(*int)(p))
}

func dumpFloat64(p unsafe.Pointer) {
	fmt.Printf("\t data: %v\n", *(*float64)(p))
}
