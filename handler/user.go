package handler

import (
	dblayer "Distributed-fileserver/db"
	"Distributed-fileserver/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
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
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	enc_passwd := util.Sha1([]byte(password + pwd_salt))


	pwdChecked := dblayer.UserSignin(username, enc_passwd)

	if pwdChecked == false{
		w.Write([]byte("FAILED"))
	}

	token := GenToken(username)
	upRes := dblayer.UpdateToken(username, token)
	if !upRes{
		w.Write([]byte("FAILED"))
		return
	}

	//w.Write([]byte("http://"+r.Host+"/static/view/home.html"))

	resp := util.RespMsg{
		Code:0,
		Msg:"ok",
		Data: struct {
			Location string
			Username string
			Token string
		}{
			Location:"http://"+r.Host+"/static/view/home.html",
			Username:username,
			Token:token,
		},
	}
	w.Write(resp.JSONBytes())
}

func GenToken(username string) string{
	//生成40位的token
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

//查询用户信息
func UserInfoHandler(w http.ResponseWriter, r *http.Request){
	r.ParseForm()
	username := r.Form.Get("username")
	token := r.Form.Get("token")

	valid := IsTokenValid(token)
	if !valid{
		w.WriteHeader(http.StatusForbidden)
		return
	}


}

//验证token有效性
func IsTokenValid(token string) bool {
	//判断token是否过期
	return true
}