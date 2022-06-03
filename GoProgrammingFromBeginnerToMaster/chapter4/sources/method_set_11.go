package main

import "fmt"

type T struct{}

func (T) M1()  {}
func (*T) M2() {}

type Interface interface {
	M1()
	M2()
}

type T1 T
type Interface1 Interface

func showDefinedType() {
	var t T
	var pt *T
	var t1 T1
	var pt1 *T1

	DumpMethodSet(&t)
	DumpMethodSet(&t1)

	DumpMethodSet(&pt)
	DumpMethodSet(&pt1)

	DumpMethodSet((*Interface)(nil))
	DumpMethodSet((*Interface1)(nil))

	var i Interface = (Interface1)(nil)
	_ = i
}

type TU struct{}

func (TU) M1() { fmt.Println("TU M1") }

type IU interface {
	M1()
	M2()
}

type T11 TU // type define

type T12 = TU // type alias

type I11 IU

func showDefinedType2() {
	var t T11
	var i I11
	var t2 T12 = T12{}
	t = T11{}
	// t.M1()
	t2.M1()

	DumpMethodSet2(&t)
	DumpMethodSet2(&i)
	DumpMethodSet2((*I11)(nil))
	DumpMethodSet2(&t2)
}

func main() {
	showDefinedType2()
}
