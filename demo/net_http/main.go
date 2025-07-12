package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GET /hello
func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	fmt.Fprintln(w, "Hello, World!")
}

// POST /echo
func echoHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	w.Write(body)
}

// GET /json
func jsonHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// レスポンス用の構造体
	type Response struct {
		Message string `json:"message"`
		Status  string `json:"status"`
	}

	resp := Response{
		Message: "Hello, JSON!",
		Status:  "success",
	}

	// JSONヘッダーをセット
	w.Header().Set("Content-Type", "application/json")

	// JSONエンコードしてレスポンス
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/echo", echoHandler)
	http.HandleFunc("/json", jsonHandler)

	fmt.Println("🚀 サーバー起動: http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
