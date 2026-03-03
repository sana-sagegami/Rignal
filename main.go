package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq" // PostgreSQLドライバー
)

func main() {
	// 1. データベース接続文字列（docker-composeの設定に合わせる）
	connStr := "host=db port=5432 user=sana password=zenpassword dbname=auto_zen_db sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 2. テーブルがなければ作成する（今回は簡易的に起動時にチェック）
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS zen_logs (id SERIAL PRIMARY KEY, timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP)")
	if err != nil {
		log.Fatal(err)
	}

	// 3. データを保存するエンドポイント
	http.HandleFunc("/save", func(w http.ResponseWriter, r *http.Request) {
		_, err := db.Exec("INSERT INTO zen_logs DEFAULT VALUES")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, `{"status": "saved", "message": "Zen mode record added! 🧘‍♀️"}`)
	})

	// 4. 起動
	fmt.Println("Server is running on http://localhost:8081 ...")
	http.ListenAndServe(":8081", nil)
}