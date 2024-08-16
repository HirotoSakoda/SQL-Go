package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load() // 環境変数の読み込み
	if err != nil {
		panic("Failed to load env file")
	}

	// 環境変数から接続情報を取得
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	// MySQLの接続情報
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, password, host, port, dbname)

	// データベースに接続
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 接続を確認
	err = db.Ping()
	if err != nil {
		fmt.Println("接続失敗")
		return
	}

	fmt.Println("接続完了")
}
