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
	
	w.Write([]byte("demo of algorithm"))
}
