package demo

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

// 例: 大量のメモリ割り当てと複雑な操作
func BenchmarkComplexWithAllocations(b *testing.B) {
    for i := 0; i < b.N; i++ {
        var sum int
        for j := 0; j < 1000; j++ {
            s := make([]int, 1000)
            for k := 0; k < 1000; k++ {
                s[k] = i + j + k
                sum += s[k]
            }
        }
    }
}

// 例: メモリ再利用と複雑な操作
func BenchmarkComplexWithoutAllocations(b *testing.B) {
    b.StopTimer()
    s := make([]int, 1000)
    b.StartTimer()
    for i := 0; i < b.N; i++ {
        var sum int
        for j := 0; j < 1000; j++ {
            for k := 0; k < 1000; k++ {
                s[k] = i + j + k
                sum += s[k]
            }
        }
    }
}
