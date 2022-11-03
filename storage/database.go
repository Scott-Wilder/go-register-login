package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"example.com/basicWebApp/types"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB

func init() {
	err := godotenv.Load(".env")

    if err != nil {
        log.Fatal("Error loading .env file")
    }	

	pswd := os.Getenv("MYSQL_PASSWORD")
	//(driver name, data source name)
	db, err = sql.Open("mysql", "root:"+pswd+"@tcp(localhost:3306)/testaccount")
	if err != nil {
		fmt.Println("error validating sql.Open arguments")
		panic(err.Error())
	}

	err = db.Ping()
	if err != nil {
		fmt.Println("error verifying connection with db.Ping")
		panic(err.Error())
	}

	fmt.Println("Successful Connection to Database!")
}

func IsUser(username string) (error){
	var user types.User
	row := db.QueryRow("SELECT * FROM user WHERE username = ?", username)
	err := row.Scan(&user.Username)
	return err
}

func PrepareInsert(username *string, lastname *string, firstname *string, hash *string) (*sql.Stmt, error) {
	var insertStmt *sql.Stmt
	insertStmt, err := db.Prepare("INSERT INTO user (username, lastname, firstname, Hash) VALUES (?, ?, ?, ?);")
	return insertStmt, err
}

func InsertUser(insertStmt *sql.Stmt, username *string, lastname *string, firstname *string, hash *string) (error){
	var result sql.Result
	//  func (s *Stmt) Exec(args ...interface{}) (Result, error)
	result, err := insertStmt.Exec(username, lastname, firstname, hash)
	rowsAff, _ := result.RowsAffected()
	lastIns, _ := result.LastInsertId()
	fmt.Println("rowsAff:", rowsAff)
	fmt.Println("lastIns:", lastIns)
	fmt.Println("err:", err)
	return err
}

func GetHash(username string) (string, error) {
	var hash string
	stmt := "SELECT Hash FROM user WHERE username = ?"
	row := db.QueryRow(stmt, username)
	err := row.Scan(&hash)
	return hash, err

}