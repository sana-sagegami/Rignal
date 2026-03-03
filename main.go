package main

import (
	"fmt"
	"net/http"
)

func main() {
	// 「/」にアクセスが来たときの返事（JSON形式）
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "success", "message": "Auto-Zen API is active! 🧘‍♀️"}`))
	})

	// 8080番ポートでサーバーを起動
	fmt.Println("Server is running on http://localhost:8081 ...")
	http.ListenAndServe(":8081", nil)
}