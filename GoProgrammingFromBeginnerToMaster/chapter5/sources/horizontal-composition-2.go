package main

import (
	"bytes"
	"io"
	"log"
	"os"
	"strings"
)

func CapReader(r io.Reader) io.Reader {
	return &capitalizedReader{r: r}
}

type capitalizedReader struct {
	r io.Reader
}

func (r *capitalizedReader) Read(p []byte) (int, error) {
	n, err := r.r.Read(p)
	if err != nil {
		return 0, err
	}

	q := bytes.ToUpper(p)
	for i, v := range q {
		p[i] = v
	}
	return n, err
}

func showCapReader() {
	r := strings.NewReader("hello, gopher!\n")
	r1 := CapReader(io.LimitReader(r, 4))
	if _, err := io.Copy(os.Stdout, r1); err != nil {
		log.Fatal(err)
	}
}

func NewCapReader(n int64) *CapReader2 {
	return &CapReader2{
		r: io.LimitReader(strings.NewReader("Hello, world"), n),
	}
}

type CapReader2 struct {
	r io.Reader
}

func (r CapReader2) Read(p []byte) (int, error) {
	n, err := r.r.Read(p)
	if err != nil {
		return 0, err
	}
	q := bytes.ToUpper(p)
	for i, v := range q {
		p[i] = v
	}
	return n, nil
}

func showCapReader2() {
	// r := NewCapReader(5)
	// lr := io.LimitReader(strings.NewReader("Hello, world"), 5)
	r := CapReader2{
		r: io.LimitReader(
			strings.NewReader("Hello, world"), 5)}
	if _, err := io.Copy(os.Stdout, r); err != nil {
		log.Fatal(err)
	}
}

func main() {
	showCapReader2()
}
