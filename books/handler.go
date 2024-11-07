package books

import (
	"fmt"
	"net/http"
)

func BooksDemoHandler(w http.ResponseWriter, r *http.Request) {
	hoge := "hoge"
	w.Write([]byte(fmt.Sprintf("hoge: %v\n", hoge)))
	w.Write([]byte("=========================\n"))
}
