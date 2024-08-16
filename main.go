package main

import (
	// "database/sql"
	"fmt"
	"log"
	"os"

	// _ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/jinzhu/gorm"
)

func main() (database *gorm.DB) {
	// MySQLの接続情報
	err := godotenv.Load(".env")
	
	// もし err がnilではないなら、コメント出力。
	if err != nil {
		fmt.Printf("読み込めませんでした: %v", err)
	}
		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",

		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		)
	
	fmt.Println(dsn)
	// データベースに接続
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 接続を確認
	err = db.Ping()
	if err != nil {
		fmt.Println("接続失敗")
		return
	} else {
		fmt.Println("接続完了")
	}
	return db
}
