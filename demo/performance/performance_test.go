package performance

import (
    "testing"
)

const size = 100000

// ベンチマーク関数: append の場合
func BenchmarkAppend(b *testing.B) {
    b.ReportAllocs()  // メモリ割り当てをレポートするように設定
    for i := 0; i < b.N; i++ {
        var slice []int
        for j := 0; j < size; j++ {
            slice = append(slice, j)
        }
    }
}

// ベンチマーク関数: hoge[i] = foo の場合
func BenchmarkIndexAssignment(b *testing.B) {
    b.ReportAllocs()  // メモリ割り当てをレポートするように設定

    for i := 0; i < b.N; i++ {
        slice := make([]int, size)
        for j := 0; j < size; j++ {
            slice[j] = j
        }
    }
}