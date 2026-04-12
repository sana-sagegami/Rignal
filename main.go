package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

// データの形を定義
type ZenRecord struct {
	ID        int    `json:"id"`
	Task      string `json:"task"`
	Duration  int    `json:"duration"`
	Timestamp string `json:"timestamp"`
}

func main() {
	// データベース接続
	connStr := "host=db port=5432 user=sana password=zenpassword dbname=auto_zen_db sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// テーブル作成（task, duration を含む）
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS zen_logs (id SERIAL PRIMARY KEY, task TEXT, duration INT, timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		log.Fatal(err)
	}

	// 1. 保存エンドポイント (POST)
	http.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POSTメソッドのみ受け付けます", http.StatusMethodNotAllowed)
			return
		}
		var record ZenRecord
		if err := json.NewDecoder(r.Body).Decode(&record); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err := db.Exec("INSERT INTO zen_logs (task, duration) VALUES ($1, $2)", record.Task, record.Duration)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, `{"status": "success", "message": "Saved: %s!"}`, record.Task)
	})

	// 2. 一覧取得エンドポイント (GET)
	http.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT id, task, duration, timestamp FROM zen_logs ORDER BY timestamp DESC")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var logs []ZenRecord
		for rows.Next() {
			var l ZenRecord
			rows.Scan(&l.ID, &l.Task, &l.Duration, &l.Timestamp)
			logs = append(logs, l)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(logs)
	})

	// 3. 削除エンドポイント (DELETE)
	http.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("id") // URLの末尾に ?id=1 などをつける
		_, err := db.Exec("DELETE FROM zen_logs WHERE id = $1", id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, `{"status": "success", "message": "Log ID %s deleted"}`, id)
	})

	fmt.Println("Server is running on http://localhost:8081 ...")
	http.ListenAndServe(":8081", nil)
}
