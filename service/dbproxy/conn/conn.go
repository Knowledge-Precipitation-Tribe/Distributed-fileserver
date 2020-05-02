package mysql

import (
	"Distributed-fileserver/service/dbproxy/customLog"
	"database/sql"
	"go.uber.org/zap"
	"os"

	cfg "Distributed-fileserver/service/dbproxy/config"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDBConn() {
	db, _ = sql.Open("mysql", cfg.MySQLSource)
	db.SetMaxOpenConns(1000)
	err := db.Ping()
	if err != nil {
		customLog.Logger.Error("InitDbConn请求失败", zap.Error(err))
		os.Exit(1)
	}
}

// DBConn : 返回数据库连接对象
func DBConn() *sql.DB {
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
		customLog.Logger.Error("dbproxy checkErr", zap.Error(err))
		panic(err)
	}
}
