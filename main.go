package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type User struct {
	ID    int
	Name  string
	Email string
}

var db *sql.DB

func connect() (*sql.DB, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("読み込めませんでした: %v", err)
	}
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func createDatabase() error {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)

	// デフォルトのデータベース名なしで接続
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("CREATE DATABASE IF NOT EXISTS mydatabase")
	if err != nil {
		return fmt.Errorf("データベース作成に失敗しました: %v", err)
	}

	log.Println("データベース 'mydatabase' を作成しました（存在しない場合）")
	return nil
}

// ユーザー一覧を表示
func listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := fetchUsers()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("ユーザー一覧の取得に失敗しました: %v", err)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/list.html"))
	err = tmpl.Execute(w, users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("ユーザー一覧のテンプレート実行に失敗しました: %v", err)
		return
	}

	log.Println("ユーザー一覧を表示しました")
}

// 新規ユーザー作成フォームの表示と作成処理
func createUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		email := r.FormValue("email")

		_, err := createUser(name, email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("ユーザーの作成に失敗しました: %v", err)
			return
		}

		log.Printf("ユーザーを作成しました: %s (%s)", name, email)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/create.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("ユーザー作成フォームのテンプレート実行に失敗しました: %v", err)
		return
	}

	log.Println("新規ユーザー作成フォームを表示しました")
}

// ユーザー編集フォームの表示と編集処理
func editUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id, _ := strconv.Atoi(r.FormValue("id"))
		name := r.FormValue("name")
		email := r.FormValue("email")

		err := updateUser(id, name, email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("ユーザーの更新に失敗しました: %v", err)
			return
		}

		log.Printf("ユーザーを更新しました: %d - %s (%s)", id, name, email)
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
	user, err := getUserByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("ユーザーの取得に失敗しました: %v", err)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/edit.html"))
	err = tmpl.Execute(w, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("ユーザー編集フォームのテンプレート実行に失敗しました: %v", err)
		return
	}

	log.Printf("ユーザー編集フォームを表示しました: %d", id)
}

// ユーザー削除
func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	err := deleteUser(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("ユーザーの削除に失敗しました: %v", err)
		return
	}

	log.Printf("ユーザーを削除しました: %d", id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// CRUD操作の関数
func fetchUsers() ([]User, error) {
	rows, err := db.Query("SELECT id, name, email FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Email)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func getUserByID(id int) (User, error) {
	var user User
	err := db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", id).Scan(&user.ID, &user.Name, &user.Email)
	if err != nil {
		return user, err
	}
	return user, nil
}

func createUser(name, email string) (int64, error) {
	result, err := db.Exec("INSERT INTO users (name, email) VALUES (?, ?)", name, email)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

func updateUser(id int, name, email string) error {
	_, err := db.Exec("UPDATE users SET name = ?, email = ? WHERE id = ?", name, email, id)
	return err
}

func deleteUser(id int) error {
	_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

func main() {
	// データベース作成
	//err := createDatabase()
	//if err != nil {
	//	log.Fatal(err)
	//}

	//データベースに接続
	var err error
	db, err = connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/", listUsers)
	http.HandleFunc("/create", createUserHandler)
	http.HandleFunc("/edit", editUserHandler)
	http.HandleFunc("/delete", deleteUserHandler)

	fmt.Println("サーバーがポート8080で起動しました")
	log.Println("サーバーがポート8080で起動しました")
	http.ListenAndServe(":8080", nil)
}
