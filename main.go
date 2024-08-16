package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func connect() {
	for i := 0; i < 10; i++ {
		err := tryConnect()
		if err == nil {
			fmt.Println("接続完了")
			return
		}
		fmt.Println("接続失敗、再試行...")
		time.Sleep(2 * time.Second)
	}
	log.Fatal("接続失敗")
}

func tryConnect() error {
	// MySQLの接続情報
	err := godotenv.Load(".env")
	
	// もし err がnilではないなら、コメント出力。
	if err != nil {
		return fmt.Errorf("読み込めませんでした: %v", err)
	}
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",

		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	// データベースに接続
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	// 接続を確認
	err = db.Ping()
	return err
}

func main() {
	connect()
}
