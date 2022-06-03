package main

import "fmt"

type T1 struct{}

func (T1) T1M1()   { println("T1's M1") }
func (T1) T1M2()   { println("T1's M2") }
func (*T1) PT1M3() { println("PT1's M3") }

type T2 struct{}

func (T2) T2M1()   { println("T2's M1") }
func (T2) T2M2()   { println("T2's M2") }
func (*T2) PT2M3() { println("PT2's M3") }

type T struct {
	T1
	*T2
}

func showStructEmbedStruct() {
	t := T{
		T1: T1{},
		T2: &T2{},
	}

	println("call method through t:")
	t.T1M1()
	t.T1M2()
	t.PT1M3()
	t.T2M1()
	t.T2M2()
	t.PT2M3()

	println("\ncall method through pt:")
	pt := &t
	pt.T1M1()
	pt.T1M2()
	pt.PT1M3()
	pt.T2M1()
	pt.T2M2()
	pt.PT2M3()
	println("")

	var t1 T1
	var pt1 *T1
	DumpMethodSet(&t1)
	DumpMethodSet(&pt1)

	var t2 T2
	var pt2 *T2
	DumpMethodSet(&t2)
	DumpMethodSet(&pt2)

	DumpMethodSet(&t)
	DumpMethodSet(&pt)

}

type Interface1 interface {
	M1()
	M2()
	M3()
}

type S1 struct {
}

func (S1) M1() {
	fmt.Println("S1 M1")
}

func (S1) M2() {
	fmt.Println("S1 M2")
}

func (*S1) M3() {
	fmt.Println("*S1 M3")
}

type Interface2 interface {
	M4()
	M5()
	M6()
}

type S2 struct {
}

func (S2) M4() {
	fmt.Println("S2 M4")
}

func (S2) M5() {
	fmt.Println("S2 M5")
}

func (*S2) M6() {
	fmt.Println("*S2 M6")
}

type Struct struct {
	S1
	*S2
}

func showStructEmbedStruct2() {
	s := Struct{
		S1: S1{},
		S2: &S2{},
	}
	pts := &Struct{
		S1: S1{},
		S2: &S2{},
	}
	s.M1()
	s.M2()
	s.M3()
	s.M4()
	s.M5()
	s.M6()
	fmt.Printf("\n")
	pts.M1()
	pts.M2()
	pts.M3()
	pts.M4()
	pts.M5()
	pts.M6()

	DumpMethodSet2(&s)
	DumpMethodSet2(&pts)
}

func main() {
	showStructEmbedStruct2()
}
