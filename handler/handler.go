package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

//处理文件上传
func UploadHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == "GET"{
		data, err := ioutil.ReadFile("./static/view/upload.html")
		if err != nil{
			io.WriteString(w, "internel server error")
			return
		}
		io.WriteString(w, string(data))
	}else{
		file, head, err := r.FormFile("file")
		if err != nil{
			fmt.Println(err)
			return
		}
		defer file.Close()

		newFile, err := os.Create("/tmp/" + head.Filename)
		if err != nil{
			fmt.Println(err)
			return
		}
		defer newFile.Close()

		_, err = io.Copy(newFile, file)
		if err != nil{
			fmt.Println(err)
			return
		}
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

//上传成功
func UploadSucHandler(w http.ResponseWriter, r * http.Request){
	io.WriteString(w, "upload success")
}