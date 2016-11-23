package util

import "testing"

func BenchmarkDenseVector_1in10(b *testing.B) {
	benchmarkVector(b, NewDenseVector(), 10, 1)
}
func BenchmarkDenseVector_10in10(b *testing.B) {
	benchmarkVector(b, NewDenseVector(), 10, 10)
}
func BenchmarkDenseVector_1in100(b *testing.B) {
	benchmarkVector(b, NewDenseVector(), 100, 1)
}
func BenchmarkDenseVector_10in100(b *testing.B) {
	benchmarkVector(b, NewDenseVector(), 100, 10)
}
func BenchmarkDenseVector_50in100(b *testing.B) {
	benchmarkVector(b, NewDenseVector(), 100, 50)
}

func BenchmarkSparseVector_1in10(b *testing.B) {
	benchmarkVector(b, NewSparseVector(), 10, 1)
}
func BenchmarkSparseVector_10in10(b *testing.B) {
	benchmarkVector(b, NewSparseVector(), 10, 10)
}
func BenchmarkSparseVector_1in100(b *testing.B) {
	benchmarkVector(b, NewSparseVector(), 100, 1)
}
func BenchmarkSparseVector_10in100(b *testing.B) {
	benchmarkVector(b, NewSparseVector(), 100, 10)
}
func BenchmarkSparseVector_50in100(b *testing.B) {
	benchmarkVector(b, NewSparseVector(), 100, 50)
}

func benchmarkVector(b *testing.B, v Vector, n, m int) {
	x := n / m
	for ; n > 0; n-- {
		if n%x == 0 {
			v = v.Set(n, 1.1)
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		v.Variance()
	}
}
