package main

import "fmt"

func odd(i int) {
	fmt.Println("odd number: ", i)
}

func even(i int) {
	fmt.Println("even number: ", i)
}

func main() {

	sl := []int{2, 34, 34, 5}
	for _, v := range sl {
		switch v {
		case 1, 3, 5, 7:
			odd(v)
		case 2, 4, 6, 8:
			even(v)
		}
	}

}
