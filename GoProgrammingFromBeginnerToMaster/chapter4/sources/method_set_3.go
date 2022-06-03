package main

import (
	"io"
)

func showInterfaceEmbeddingInterface() {
	DumpMethodSet((*io.Writer)(nil))
	DumpMethodSet((*io.Reader)(nil))
	DumpMethodSet((*io.Closer)(nil))
	DumpMethodSet((*io.ReadWriter)(nil))
	DumpMethodSet((*io.ReadWriteCloser)(nil))
}

type MyInt int

func showInterfaceEmbeddingInterface2() {
	// i := 1
	// i2 := MyInt(i)

	DumpMethodSet2((*io.Reader)(nil))
	DumpMethodSet2((*io.Writer)(nil))
	DumpMethodSet2((*io.Closer)(nil))
	DumpMethodSet2((*io.ReadWriteCloser)(nil))
}

func main() {
	showInterfaceEmbeddingInterface2()
}
