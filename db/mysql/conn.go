package db

import(
	"Distributed-fileserver/config"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)
var db *sql.DB

func init()  {
	datasource := config.MySQLUser + ":" + config.MySQLPassword + "@tcp(" + config.MySQLHost + ")/" + config.DBName + "?charset="  + config.DBCharset
	db, _ = sql.Open("mysql", datasource)
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil{
		fmt.Println(err.Error())
		fmt.Println("aaaaaa")
		os.Exit(1)
	}
}

func DBConn() *sql.DB{
	return db
}

func ParseRows(rows *sql.Rows) []map[string]interface{} {
	columns, _ := rows.Columns()
	scanArgs := make([]interface{}, len(columns))
	values := make([]interface{}, len(columns))
	for j := range values {
		scanArgs[j] = &values[j]
	}

	record := make(map[string]interface{})
	records := make([]map[string]interface{}, 0)
	for rows.Next() {
		//将行数据保存到record字典
		err := rows.Scan(scanArgs...)
		checkErr(err)

		for i, col := range values {
			if col != nil {
				record[columns[i]] = col
			}
		}
		records = append(records, record)
	}
	return records
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}
