package main

import (
	"fmt"
	"strings"
	// "golang.org/x/text/cases"
)

func concat(sep string, args ...interface{}) string {
	var result string
	for i, v := range args {
		if i != 0 {
			result += sep
		}
		switch v.(type) {
		case int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64:
			result += fmt.Sprintf("%d", v)
		case string:
			result += fmt.Sprintf("%s", v)
		case []int:
			ints := v.([]int)
			for i, v := range ints {
				if i != 0 {
					result += sep
				}
				result += fmt.Sprintf("%d", v)
			}
		case []string:
			strs := v.([]string)
			result += strings.Join(strs, sep)
		default:
			fmt.Printf("the argument type [%T] is not supported", v)
			return ""
		}
	}
	return result
}

func concat2(sep string, args ...interface{}) string {
	res := ""
	for i, v := range args {
		if i != 0 {
			res += sep
		}
		switch v.(type) {
		case int, int16, int32, int64, uint, uint16, uint32, uint64:
			res += fmt.Sprintf("%d", v)
		case string:
			res += fmt.Sprintf("%s", v)
		case []int:
			ints := v.([]int)
			for _, val := range ints {
				res += fmt.Sprintf("%d", val)
			}
		case []string:
			strs := v.([]string)
			for _, val := range strs {
				res += val
			}
		case []interface{}:
			is := v.([]interface{})
			for _, val := range is {
				v1, ok := val.(string)
				if ok {
					res += v1
				}
			}
		default:
			fmt.Println("type not match")
		}
	}
	return res
}

func showDefaultAndCustomVariable() {
	println(concat("-", 1, 2))
	println(concat("-", "hello", "gopher"))
	println(concat("-", "hello", 1, uint32(2),
		[]int{11, 12, 13}, 17,
		[]string{"robot", "ai", "ml"},
		"hacker", 33))
}

func showDefaultAndCustomVariable2() {
	println(concat2("-", 1, 2, struct{ i int }{i: 1}))
	println(concat2("-", 1, 2))
	println(concat2("-", "hello", "gopher"))
	println(concat2("-", "hello", 1, uint32(2),
		[]int{11, 12, 13}, 17,
		[]string{"robot", "ai", "ml"},
		"hacker", 33))
}

func main() {
	// showDefaultAndCustomVariable()
	showDefaultAndCustomVariable2()
}
