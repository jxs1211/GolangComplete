package pkg3

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
	fmt.Printf("pkg3 init %s\n", s)
	return s
}

func checkInit() string {
	if c1 != "" {
		fmt.Println("pkg3 init c1")
	}
	if c2 != "" {
		fmt.Println("pkg3 init c2")
	}
	return ""
}

func init() {
	fmt.Println("pkg3 init")
}

func main() {
	fmt.Println("pkg3 main")
}

func innerFunc() {
	fmt.Println("innerFunc in the same package")
}

func ExportFunc() {
	Export()
	inner()
	fmt.Println("pkg3 file pkg3 exportFunc")
}
