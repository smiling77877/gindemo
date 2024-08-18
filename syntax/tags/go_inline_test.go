package tags

import "testing"

//go:inline
func SumV1(a, b int) int {
	return a + b
}

//go:noinline
func SumV2(a, b int) int {
	return a + b
}

var Result int

func Benchmark_TestSumInline(b *testing.B) {
	for i := 0; i < b.N; i++ {
		res := SumV1(i, i)
		Result = res
	}
}

func Benchmark_TestSumNoInline(b *testing.B) {
	for i := 0; i < b.N; i++ {
		res := SumV2(i, i)
		Result = res
	}
}
