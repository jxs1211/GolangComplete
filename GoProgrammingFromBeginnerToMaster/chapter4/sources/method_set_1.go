package main

type Interface interface {
	M1()
	M2()
}

type T struct{}

func (t T) M1()  {}
func (t *T) M2() {}

func showInterfaceTypeStruct() {
	var t T
	var pt *T
	var i Interface

	i = t
	i = pt
}

type Interface2 interface {
	M1()
	M2()
}

type T2 struct {
}

func (t T2) M1() {

}

func (t *T2) M2() {

}

func showInterfaceTypeStruct2() {
	var t T2
	var pt *T2
	var i Interface2

	i = t
	i = pt
}

func main() {
	showInterfaceTypeStruct2()
}
