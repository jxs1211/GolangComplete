package main

import "fmt"

func arrayRangeExpression() {
	var a = [5]int{1, 2, 3, 4, 5}
	var r [5]int

	fmt.Println("arrayRangeExpression result:")
	fmt.Println("a = ", a)

	for i, v := range a {
		if i == 0 {
			a[1] = 12
			a[2] = 13
		}

		r[i] = v
	}

	fmt.Println("r = ", r)
	fmt.Println("a = ", a)
}

func pointerToArrayRangeExpression() {
	var a = [5]int{1, 2, 3, 4, 5}
	var r [5]int

	fmt.Println("pointerToArrayRangeExpression result:")
	fmt.Println("a = ", a)

	for i, v := range &a {
		if i == 0 {
			a[1] = 12
			a[2] = 13
		}

		r[i] = v
	}

	fmt.Println("r = ", r)
	fmt.Println("a = ", a)
}

func sliceRangeExpression() {
	var a = [5]int{1, 2, 3, 4, 5}
	var r [5]int

	fmt.Println("sliceRangeExpression result:")
	fmt.Println("a = ", a)

	for i, v := range a[:] {
		if i == 0 {
			a[1] = 12
			a[2] = 13
		}

		r[i] = v
	}

	fmt.Println("r = ", r)
	fmt.Println("a = ", a)
}

func sliceLenChangeRangeExpression() {
	var a = []int{1, 2, 3, 4, 5}
	var r = make([]int, 0)

	fmt.Println("sliceLenChangeRangeExpression result:")
	fmt.Println("a = ", a)

	for i, v := range a {
		if i == 0 {
			a = append(a, 6, 7)
		}

		r = append(r, v)
	}

	fmt.Println("r = ", r)
	fmt.Println("a = ", a)
}

func arrayRangeExpression2() {
	a := [4]int{1, 2, 3, 4}
	fmt.Println("a: ", a)
	var r [4]int
	for i, v := range a {
		if i == 0 {
			a[1] = 22
			a[2] = 33
		}
		r[i] = v
	}
	fmt.Println("r: ", r)
	fmt.Println("a: ", a)
}

func pointerToArrayRangeExpression2() {
	a := [4]int{1, 2, 3, 4}
	fmt.Println("a: ", a)
	var r [4]int
	for i, v := range &a {
		if i == 0 {
			a[1] = 22
			a[2] = 33
		}
		r[i] = v
	}
	fmt.Println("r: ", r)
	fmt.Println("a: ", a)
}

func sliceRangeExpression2() {
	a := []int{1, 2, 3, 4}
	fmt.Println("a: ", a)
	var r []int
	for i, v := range a {
		if i == 0 {
			a[1] = 22
			a[2] = 33
		}
		r = append(r, v)
	}
	fmt.Println("r: ", r)
	fmt.Println("a: ", a)
}

func sliceLenChangeRangeExpression2() {
	a := []int{1, 2, 3, 4}
	fmt.Println("a: ", a)
	var r []int
	for i, v := range a { // the a' underlying strut's length is 4 when copying
		if i == 0 {
			a = append(a, 5)
			a = append(a, 6)
		}
		r = append(r, v)
	}
	fmt.Println("r: ", r)
	fmt.Println("a: ", a)
}

func main() {
	arrayRangeExpression2()
	pointerToArrayRangeExpression2()
	sliceRangeExpression2()
	sliceLenChangeRangeExpression2()
	// arrayRangeExpression()
	// pointerToArrayRangeExpression()
	// sliceRangeExpression()
	// sliceLenChangeRangeExpression()
}
