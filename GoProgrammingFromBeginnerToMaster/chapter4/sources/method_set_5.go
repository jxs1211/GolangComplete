package main

type Interface interface {
	M1()
	M2()
}

type T struct {
	Interface
}

func (T) M3() {}

func (*T) M4() {}

func main() {
	DumpMethodSet2((*Interface)(nil))
	var t T
	var pt *T
	DumpMethodSet2(&t)
	DumpMethodSet2(&pt)
}
