package db

import (
	mydb "Distributed-fileserver/db/mysql"
	"fmt"
	"time"
)

type UserFile struct {
	UserName string
	FileHash string
	FileName string
	FileSize int64
	UploadAt string
	LastUpdated string
}

//更新用户文件表
func OnUserFileUploadFinished(username string, filehash string,
	filename string, filesize int64) bool{
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_user_file (`user_name`, `file_sha1`, `file_name`," +
			"`file_size`, `upload_at`) values (?,?,?,?,?)")
	if err != nil{
		fmt.Println(err.Error())
		return false
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, filehash, filename, filesize, time.Now())
	if err != nil{
		fmt.Println(err.Error())
		return false
	}
	return true
}