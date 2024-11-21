package main

import "testing"

const size = 10000000

// ベンチマーク関数: append の場合
func BenchmarkAppend(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var slice []int
		for j := 0; j < size; j++ {
			slice = append(slice, j)
		}
	}
}

// ベンチマーク関数: slice[i] = int の場合
func BenchmarkIndexAssignment(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
        slice := make([]int, size)
        for j := 0; j < size; j++ {
            slice[j] = j
        }
    }
}
