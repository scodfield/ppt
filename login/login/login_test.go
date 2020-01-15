package main

import (
	"testing"
)

func Benchmark_Login(b *testing.B) {
	var n int 
	for i := 0; i < b.N; i++ {
		n++
	}
}
