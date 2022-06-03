package main

type Interface interface {
	M1()
	M2()
}

type T struct{}

func (t T) M1()  {}
func (t *T) M2() {}

func main() {
	var t T
	var pt *T
	DumpMethodSet2(&t)
	DumpMethodSet2(&pt)
	DumpMethodSet2((*Interface)(nil))
}
