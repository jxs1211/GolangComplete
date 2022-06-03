package main

import "fmt"

const (
	PI   = 3.1415
	PI_2 = 3.1415 / (2 * iota)
	PI_4
)

const (
	a = 0
	b = iota
	c
)

func ShowFloatTypeEnums() {
	fmt.Println(PI)
	fmt.Println(PI_2)
	fmt.Println(PI_4)
	fmt.Println(a, b, c)
}

type Person struct {
	name string
	age  int
}

func showCompositeType() {
	m := make(map[int]string)
	b := [3]byte{'s', 'h', 'e'}
	b2 := []byte{'s'}
	p := Person{
		name: "shen",
		age:  20,
	}
	fmt.Println(m, b, b2, p)

	b3 := []Person{
		{"shen", 20},
		{"xian", 21},
	}
	b4 := [][]int{
		{1, 2, 3},
		{4, 5, 6},
	}
	b5 := []map[int]Person{
		{1: {"shen", 20}},
		{2: {"xian", 21}},
	}
	fmt.Printf("%v\n", b3)
	fmt.Println(b4)
	fmt.Println(b5)
}

func showSliceLenAndCap() {
	s := make([]int, 0, 4)
	s2 := make([]int, 4)
	s3 := []int{1, 2, 3}
	fmt.Printf("s: %v len: %d, cap: %d\n", s, len(s), cap(s))
	fmt.Printf("s2: %v len: %d, cap: %d\n", s2, len(s2), cap(s2))
	fmt.Printf("s3: %v len: %d, cap: %d\n", s3, len(s3), cap(s3))
}

func main() {

	showSliceLenAndCap()
	// showCompositeType()
	// ShowFloatTypeEnums()
}
