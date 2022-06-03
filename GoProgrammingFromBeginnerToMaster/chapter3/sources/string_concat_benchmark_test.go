package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

var sl []string = []string{
	"Rob Pike ",
	"Robert Griesemer ",
	"Ken Thompson ",
}

func concatStringByOperator(sl []string) string {
	var s string
	for _, v := range sl {
		s += v
	}
	return s
}

func concatStringBySprintf(sl []string) string {
	var s string
	for _, v := range sl {
		s = fmt.Sprintf("%s%s", s, v)
	}
	return s
}

func concatStringByJoin(sl []string) string {
	return strings.Join(sl, "")
}

func concatStringByStringsBuilder(sl []string) string {
	var b strings.Builder
	for _, v := range sl {
		b.WriteString(v)
	}
	return b.String()
}

func concatStringByStringsBuilderWithInitSize(sl []string) string {
	var b strings.Builder
	b.Grow(64)
	for _, v := range sl {
		b.WriteString(v)
	}
	return b.String()
}

func concatStringByBytesBuffer(sl []string) string {
	var b bytes.Buffer
	for _, v := range sl {
		b.WriteString(v)
	}
	return b.String()
}

func concatStringByBytesBufferWithInitSize(sl []string) string {
	buf := make([]byte, 0, 64)
	b := bytes.NewBuffer(buf)
	for _, v := range sl {
		b.WriteString(v)
	}
	return b.String()
}

func BenchmarkConcatStringByOperator(b *testing.B) {
	for n := 0; n < b.N; n++ {
		// concatStringByOperator(sl)
		concatStringWithOperator(s)
	}
}

func BenchmarkConcatStringBySprintf(b *testing.B) {
	for n := 0; n < b.N; n++ {
		// concatStringBySprintf(sl)
		concatStringWithFormat(s)
	}
}

func BenchmarkConcatStringByJoin(b *testing.B) {
	for n := 0; n < b.N; n++ {
		// concatStringByJoin(sl)
		concatStringWithJoin(s)
	}
}

func BenchmarkConcatStringByStringsBuilder(b *testing.B) {
	for n := 0; n < b.N; n++ {
		// concatStringByStringsBuilder(sl)
		concatStringWithStringBuilder(s)
	}
}

func BenchmarkConcatStringByStringsBuilderWithInitSize(b *testing.B) {
	for n := 0; n < b.N; n++ {
		// concatStringByStringsBuilderWithInitSize(sl)
		concatStringWithStringBuilderInitSize(s)
	}
}

func BenchmarkConcatStringByBytesBuffer(b *testing.B) {
	for n := 0; n < b.N; n++ {
		// concatStringByBytesBuffer(sl)
		concatStringWithBytesBuffer(s)
	}
}

func BenchmarkConcatStringByBytesBufferWithInitSize(b *testing.B) {
	for n := 0; n < b.N; n++ {
		// concatStringByBytesBufferWithInitSize(sl)
		concatStringWithBytesBufferInitSize(s)
	}
}

var s []string = []string{"shen沈", "xian先", "jie捷"}

func BenchmarkconcatStringWithOperator(b *testing.B) {
	for i := 0; i < b.N; i++ {
		concatStringWithOperator(s)
	}
}

func concatStringWithOperator(s []string) string {
	var res string
	for _, v := range s {
		res += v
	}
	return res
}

func concatStringWithFormat(s []string) string {
	var res string
	for _, v := range s {
		res += fmt.Sprintf("%s", v)
	}
	return res
}

func concatStringWithStringBuilder(s []string) string {
	var b strings.Builder
	for _, v := range s {
		b.WriteString(v)
	}
	return b.String()
}

func concatStringWithStringBuilderInitSize(s []string) string {
	var b strings.Builder
	b.Grow(64)
	for _, v := range s {
		b.WriteString(v)
	}
	return b.String()
}

func concatStringWithBytesBuffer(s []string) string {
	var b bytes.Buffer
	for _, v := range s {
		b.WriteString(v)
	}
	return b.String()
}

func concatStringWithBytesBufferInitSize(s []string) string {
	buf := make([]byte, 0, 64)
	b := bytes.NewBuffer(buf)
	for _, v := range s {
		b.WriteString(v)
	}
	return b.String()
}

func concatStringWithJoin(s []string) string {
	return strings.Join(sl, "")
}
