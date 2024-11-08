package handler

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/keito-isurugi/go-demo/demo/algorithm"
)

func AlgorithmDemoHandler(w http.ResponseWriter, r *http.Request) {
	// 線形探索
	lss := generateRandomArray(1, 5)
	lst := generateRandomInt(1, 5)
	w.Write([]byte(fmt.Sprintf("Array: %v\n", lss)))
	w.Write([]byte(fmt.Sprintf("Search target: %d\n", lst)))

	lsr := algorithm.LinearSearch(lss, lst)
	if lsr != -1 {
		w.Write([]byte(fmt.Sprintf("二分探索(配列内の値の重複なし)：ターゲット %d はインデックス %d にあります\n", lst, lsr)))
	} else {
		w.Write([]byte("二分探索(配列内の値の重複なし)：該当する値は存在しません。"))
	}
	w.Write([]byte("=========================\n"))


	// バブルソート(昇順)
	bass := generateRandomArray(1, 50)
	w.Write([]byte(fmt.Sprintf("Array: %v\n", bass)))
	bssr := algorithm.BubbleAscSort(bass)
	w.Write([]byte(fmt.Sprintf("バブルソート(昇順): %d\n", bssr)))	
	w.Write([]byte("=========================\n"))


	// バブルソート(降順)
	bdss := generateRandomArray(1, 50)
	w.Write([]byte(fmt.Sprintf("Array: %v\n", bdss)))
	bdsr := algorithm.BubbleDescSort(bdss)
	w.Write([]byte(fmt.Sprintf("バブルソート(降順): %d\n", bdsr)))
	w.Write([]byte("=========================\n"))


	// 二分探索(配列内の値の重複なし)
	bsAry := []int{1, 2, 3, 4, 5, 6}
	bsTarget := generateRandomInt(1, 5)
	w.Write([]byte(fmt.Sprintf("Array: %v\n", bsAry)))
	w.Write([]byte(fmt.Sprintf("Search target: %d\n", bsTarget)))

	bsRes := algorithm.BinarySearch(bsAry, bsTarget)
	if bsRes != -1 {
		w.Write([]byte(fmt.Sprintf("二分探索(配列内の値の重複なし)：ターゲット %d はインデックス %d にあります\n", bsTarget, bsRes)))
	} else {
		w.Write([]byte("二分探索(配列内の値の重複なし)：該当する値は存在しません。"))
	}
	w.Write([]byte("=========================\n"))
}

func generateRandomArray(min, max int) []int {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)

	randomArray := make([]int, max)
	for i := 0; i < max; i++ {
		randomArray[i] = r.Intn(max - min + 1) + min
	}
	return randomArray
}

func generateRandomInt(min, max int) int {
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	return r.Intn(max - min + 1) + min
}