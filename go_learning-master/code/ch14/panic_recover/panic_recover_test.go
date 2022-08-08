package panic_recover

import (
	"errors"
	"fmt"
	"os"
	"testing"
)

func TestPanicVxExit(t *testing.T) {

	defer func() {
		if err := recover(); err != nil {
			fmt.Println("recovered from ", err)
		}
	}()
	fmt.Println("Start")
	panic(errors.New("Something wrong!"))
	//os.Exit(-1)
	//fmt.Println("End")
}

func TestPanicVxExit2(t *testing.T) {
	defer func() {
		if e := recover(); e != nil {
			fmt.Printf("recover from error: %v\n", e)
		}
	}()
	panic(errors.New("Ops! something wrong"))
}

func TestPanicVxExit3(t *testing.T) {
	defer println("Hi")
	println("before exit 1")
	os.Exit(1) // won't call defer function
	println("after exit 1")
}
