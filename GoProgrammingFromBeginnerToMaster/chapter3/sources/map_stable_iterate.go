package main

import "fmt"

func doIteration(sl []int, m map[int]int) {
	fmt.Printf("{ ")
	for _, k := range sl {
		v, ok := m[k]
		if !ok {
			continue
		}
		fmt.Printf("[%d, %d] ", k, v)
	}
	fmt.Printf("}\n")
}

func showMapStableRead() {
	var sl []int
	m := map[int]int{
		1: 11,
		2: 12,
		3: 13,
	}

	for k, _ := range m {
		sl = append(sl, k)
	}
	fmt.Println(sl)
	for i := 0; i < 3; i++ {
		doIteration(sl, m)
	}
}

func showMapStableRead2() {
	m := map[int]int{
		1: 1,
		2: 2,
		3: 3,
	}
	s := make([]int, 0, len(m))
	for k, _ := range m {
		s = append(s, k)
	}
	fmt.Println(s)
	for i := 0; i < 3; i++ {
		fmt.Print("[")
		for i := range s {
			fmt.Printf(" %d, %d ", s[i], m[s[i]])
		}
		fmt.Println("]")
	}
}

func main() {
	// showMapStableRead2()
	showMapStableRead()
}
