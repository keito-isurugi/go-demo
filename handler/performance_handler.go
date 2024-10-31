package handler

import (
	"net/http"

	"github.com/keito-isurugi/go-demo/demo/performance"
)

func PerformanceDemoHandler(w http.ResponseWriter, r *http.Request) {
	performance.PerformanceDemo()
	w.Write([]byte("demo of performance"))
}
