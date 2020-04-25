package db

import (
	mydb "Distributed-fileserver/db/mysql"
	"fmt"
)

// 用户注册，通过用户名和密码
func UserSignup(username string, password string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_user (`user_name`, `user_pwd`) " +
			"values (?,?)")
	if err != nil{
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	result, err := stmt.Exec(username, password)
	if err != nil{
		fmt.Println(err.Error())
		return false
	}

	if rowsAffected, err := result.RowsAffected(); nil == err && rowsAffected > 0{
		return true
	}
	return false

}