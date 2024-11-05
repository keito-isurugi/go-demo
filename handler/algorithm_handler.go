package handler

import (
	"fmt"
	"net/http"

	"github.com/keito-isurugi/go-demo/demo/algorithm"
)

func AlgorithmDemoHandler(w http.ResponseWriter, r *http.Request) {
	lss := []int{ 1, 2, 3, 4, 5 }
	lst := 4
	lsr, err := algorithm.LinearSearch(lss, lst)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(lsr)
	fmt.Println("=========================")

	bass := []int{ 8, 4, 5, 2, 9, 10}
	bssr := algorithm.BubbleAscSort(bass)
	fmt.Println(bssr)
	fmt.Println("=========================")
	
	bdss := []int{ 8, 4, 5, 2, 9, 10}
	bdsr := algorithm.BubbleDescSort(bdss)
	fmt.Println(bdsr)
	fmt.Println("=========================")

	w.Write([]byte("demo of algorithm"))
}
