package db

import(
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"os"
)
var db *sql.DB

func init()  {
	db, _ = sql.Open("mysql", "root:root@tcp(139.9.131.190:3306)/fileserver?charset=utf8")
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil{
		fmt.Println(err)
		os.Exit(1)
	}
}

func DBConn() *sql.DB{
	return db
}