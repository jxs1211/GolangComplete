package a

import (
	"test/b"
	c "test/b/internal"
	"test/internal/pkg/print"
)

func APrintln(s string) {
	print.Println(s)
	b.BPrintln(s)
	c.CPrintln(s) // internal limit
}
