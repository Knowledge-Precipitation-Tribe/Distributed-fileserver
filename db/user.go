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

//判断密码是否一致
func UserSignin(username string, encpwd string) bool{
	stmt, err := mydb.DBConn().Prepare(
		"select `user_pwd` from tbl_user where user_name = ? limit 1")
	if err != nil{
		fmt.Println(err.Error())
		return false
	}

	var pwd string
	err = stmt.QueryRow(username).Scan(&pwd)
	if err != nil{
		fmt.Println(err.Error())
		return false
	}else if pwd == ""{
		fmt.Println("username not found", username)
		return false
	}else if len(pwd) > 0 && pwd == encpwd{
		return true
	}

	//result, err := stmt.Exec(username)

	//pRows := mydb.ParseRows(result)
	//if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd{
	//	return true
	//}
	return false
}

//更新token
func UpdateToken(username string, token string) bool{
	stmt, err := mydb.DBConn().Prepare(
		"replace into tbl_user_token(`user_name`, `user_token`) values (?,?)")
	if err != nil{
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, token)
	if err != nil{
		fmt.Println(err.Error())
		return false
	}
	return true

}

type User struct {
	Username string
	Email string
	Phone string
	SignupAt string
	LastActiveAt string
	Status int
}

func GetUserInfo(username string)(User, error){
	user := User{}

	stmt, err := mydb.DBConn().Prepare(
		"select user_name, signup_at from tbl_user where user_name = ? limit = 1")
	if err != nil{
		fmt.Println(err.Error())
		return user, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(username).Scan(&user.Username, &user.SignupAt)
	if err != nil{
		return user, err
	}
	return user, nil
}