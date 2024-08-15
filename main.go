package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// MySQLの接続情報
	dsn := "mysql-container:my-secret-pw@tcp(128d8f4698f0:3306)/"

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
