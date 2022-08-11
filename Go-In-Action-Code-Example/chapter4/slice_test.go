package slice

import (
	"testing"
)

func TestSlice(t *testing.T) {
	var s []int
	for i := range s {
		println(i)
	}
	s2 := make([]int, 0)
	for i := range s2 {
		println(i)
	}
}

type S struct {
	name string
	age  *int
}

func (s *S) change() {
	s.name = "xian"
	a := 100
	s.age = &a
}

func TestNoprimitive(t *testing.T) {
	a := 10
	s := S{
		name: "shen",
		age:  &a,
	}
	s.change()
	t.Log(s.name, *s.age)
}
