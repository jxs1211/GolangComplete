package main

import (
	"bytes"
	"sync"
	"testing"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func writeBufFromPool(data string) {
	b := bufPool.Get().(*bytes.Buffer)
	b.Reset()
	b.WriteString(data)
	bufPool.Put(b)
}

func writeBufFromNew(data string) *bytes.Buffer {
	b := new(bytes.Buffer)
	b.WriteString(data)
	return b
}

func BenchmarkWithoutPool(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		writeBufFromNew("hello")
	}
}

func BenchmarkWithPool(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		writeBufFromPool("hello")
	}
}

var bPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

func writeDataFromPool(s string) {
	b := bPool.Get().(*bytes.Buffer)
	b.WriteString(s)
	bPool.Put(b)
}

func writeDataFromNew(s string) {
	var b bytes.Buffer
	b.WriteString(s)
}

func BenchmarkWriteDataFromPool(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(b *testing.PB) {
		for b.Next() {
			writeDataFromPool("hello")
		}
	})
}

func BenchmarkWriteDataFromNew(b *testing.B) {
	b.ReportAllocs()
	b.RunParallel(func(b *testing.PB) {
		for b.Next() {
			writeDataFromNew("hello")
		}
	})
}

func BenchmarkWriteDataFromPool2(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		writeDataFromPool("hello")
	}
}

func BenchmarkWriteDataFromNew2(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		writeDataFromNew("hello")
	}
}
