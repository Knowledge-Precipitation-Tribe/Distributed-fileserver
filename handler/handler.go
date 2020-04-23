package handler

import (
	"Distributed-fileserver/meta"
	"Distributed-fileserver/util"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
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

		fileMeta := meta.FileMeta{
			FileName:head.Filename,
			Location:"/tmp/"+head.Filename,
			UploadAt:time.Now().Format("2006-01-02 15:04:05"),
		}

		newFile, err := os.Create(fileMeta.Location)
		if err != nil{
			fmt.Println(err)
			return
		}
		defer newFile.Close()

		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil{
			fmt.Println(err)
			return
		}

		newFile.Seek(0,0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		meta.UploadFileMeta(fileMeta)

		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

//上传成功
func UploadSucHandler(w http.ResponseWriter, r * http.Request){
	io.WriteString(w, "upload success")
}

//获取文件元信息
func GetFileMetaHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()

	filehash := r.Form["filehash"][0]
	fMeta := meta.GetFileMeta(filehash)
	data, err := json.Marshal(fMeta)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

//下载文件
func DownloadHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	fsha1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fsha1)

	f, err := os.Open(fm.Location)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil{
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Description", "attachment;filename=\""+fm.FileName+"\"")
	w.Write(data)

}