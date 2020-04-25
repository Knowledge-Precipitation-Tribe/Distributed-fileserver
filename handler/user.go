package handler

import (
	dblayer "Distributed-fileserver/db"
	"Distributed-fileserver/util"
	"io/ioutil"
	"net/http"
)

const (
	pwd_salt = "*#890"
)


//处理用户注册请求
func SignupHandler(w http.ResponseWriter, r *http.Request){
	if r.Method == http.MethodGet{
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}
	r.ParseForm()

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	if len(username) < 3 || len(password) < 5{
		w.Write([]byte("invalid parameter"))
		return
	}

	enc_passwd := util.Sha1([]byte(password + pwd_salt))
	suc := dblayer.UserSignup(username, enc_passwd)
	if suc {
		w.Write([]byte("SUCCESS"))
	}else{
		w.Write([]byte("FAILED"))
	}
}


//用户登录接口
func SignInHandler(w http.ResponseWriter, r *http.Request){

}