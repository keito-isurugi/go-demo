package books

import (
	"net/http"
)

func BooksDemoHandler(w http.ResponseWriter, r *http.Request) {
	hoge := "hoge"
	if _, err := w.Write([]byte("hoge: " + hoge + "\n")); err != nil {
		return
	}
	if _, err := w.Write([]byte("=========================\n")); err != nil {
		return
	}
}
