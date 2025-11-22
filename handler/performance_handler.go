package handler

import (
	"net/http"

	"github.com/keito-isurugi/go-demo/demo"
)

func PerformanceHandler(w http.ResponseWriter, r *http.Request) {
	demo.PerformanceDemo()
	w.Write([]byte("performance comparison!"))
} 

func PerformanceProfHandler(w http.ResponseWriter, r *http.Request) {
	demo.PerformanceProfDemo()
	w.Write([]byte("performance comparison!"))
} 

func PerformanceTraceHandler(w http.ResponseWriter, r *http.Request) {
	demo.PerformanceTraceDemo()
	w.Write([]byte("performance comparison!"))
} 
