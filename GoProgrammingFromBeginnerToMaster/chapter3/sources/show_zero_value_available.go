package main

import (
	"bytes"
	"fmt"
	"net"
)

func ShowZeroSlice() {
	// var i []byte
	b := []byte("Effective Go")
	i := make([]byte, 0, len(b))
	copy(i[0:], b)
	fmt.Println(string(i), len(i), cap(i))
}

func ShowZeroSlice2() {
	var b []byte
	b = append(b, 'E')
	b = append(b, 'f')
	fmt.Println(string(b))
}

func ShowCallMethodThroughNilReceiver() {
	var a *net.TCPAddr
	fmt.Println(a)
}

func ShowZeroBytesBuffer() {
	var b bytes.Buffer
	b.Write([]byte("Effective Go"))
	fmt.Println(b.String())
}

func main() {
	ShowCallMethodThroughNilReceiver()
	// ShowZeroSlice2()
	// ShowZeroSlice()
	// ShowZeroBytesBuffer()
}
