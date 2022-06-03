package pkg2

import "fmt"

var (
	_  = checkInit()
	V1 = initVariable("v1")
	v2 = initVariable("v2")
)

const (
	c1 = "c1"
	c2 = "c2"
)

func initVariable(s string) string {
	fmt.Printf("pkg2 init %s\n", s)
	return s
}

func checkInit() string {
	if c1 != "" {
		fmt.Println("pkg2 init c1")
	}
	if c2 != "" {
		fmt.Println("pkg2 init c2")
	}
	return ""
}

func init() {
	fmt.Println("pkg2 init")
}

func main() {
	fmt.Println("pkg2 main")
}
