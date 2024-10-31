package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/keito-isurugi/go-demo/demo/time"
)

func TimeDemoHandler(w http.ResponseWriter, r *http.Request) {
	t := time.Date(2023, 12, 31, 0, 0, 0, 0, time.UTC)
	addTwoMonth := demo.AddMonthsPreservingEndOfMonth(t, 2)
	fmt.Println(t)
	fmt.Println(addTwoMonth)

	w.Write([]byte("demo of time"))
}
