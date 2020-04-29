package handler

import (
	rPool "Distributed-fileserver/cache/redis"
	dblayer "Distributed-fileserver/db"
	"Distributed-fileserver/util"
	"fmt"
	"github.com/garyburd/redigo/redis"
	"math"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

// MultipartUploadInfo: 分块上传初始化信息
type MultipartUploadInfo struct {
	FileHash   string
	FileSize   int
	UploadID   string
	ChunkSize  int
	ChunkCount int
}

const CHUNK_SIZE = 5 * 1024 * 1024

//初始化分块上传
func InitialMultipartUploadHandler(w http.ResponseWriter, r *http.Request) {
	//1 解析用户请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	fileseize, err := strconv.Atoi(r.Form.Get("filesize"))
	if err != nil {
		w.Write(util.NewRespMsg(-1, "params invalid", nil).JSONBytes())
		return
	}
	//2 获得redis的一个连接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()
	//3 获得分块上传的初始化信息
	upInfo := MultipartUploadInfo{
		FileHash:   filehash,
		FileSize:   fileseize,
		UploadID:   username + fmt.Sprintf("%x", time.Now().UnixNano()),
		ChunkSize:  CHUNK_SIZE, //每块5M
		ChunkCount: int(math.Ceil(float64(fileseize) / CHUNK_SIZE)),
	}
	//4 将初始化信息写入到redis缓存
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "chunkcount", upInfo.ChunkCount)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filehash", upInfo.FileHash)
	rConn.Do("HSET", "MP_"+upInfo.UploadID, "filesize", upInfo.FileSize)
	//5 将相应初始化数据返回到客户端
	w.Write(util.NewRespMsg(0, "OK", upInfo).JSONBytes())
}

//上传文件分块接口
func UploadPartHandler(w http.ResponseWriter, r *http.Request) {
	//1 解析用户请求参数
	r.ParseForm()
	//username := r.Form.Get("username")
	uploadID := r.Form.Get("uploadid")
	chunkIndex := r.Form.Get("index")
	//2 获得redis连接池中的一个连接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()
	//3 获得文件句柄，用户存储分块内容
	fpath := "/data/" + uploadID + "/" + chunkIndex
	os.MkdirAll(path.Dir(fpath), 0744)
	fd, err := os.Create(fpath)
	if err != nil {
		w.Write(util.NewRespMsg(-1, "upload part failed", nil).JSONBytes())
		return
	}
	defer fd.Close()

	buf := make([]byte, 1024)
	for {
		n, err := r.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil {
			break
		}
	}

	//4 更新redis缓存状态
	rConn.Do("HSET", "MP_"+uploadID, "chkidx_"+chunkIndex, 1)
	//5 返回处理结果到客户端
	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}

func CompleteUploadHandler(w http.ResponseWriter, r *http.Request) {
	//1 解析请求参数
	r.ParseForm()
	upid := r.Form.Get("uploadid")
	username := r.Form.Get("username")
	filehash := r.Form.Get("filehash")
	filesize := r.Form.Get("filesize")
	filename := r.Form.Get("filename")
	//2 获得redis连接池中的一个连接
	rConn := rPool.RedisPool().Get()
	defer rConn.Close()
	//3 通过uploadid查询redis并判断是否所有分块都已经上传完成
	data, err := redis.Values(rConn.Do("HGETALL", "MP_"+upid))
	if err != nil {
		w.Write(util.NewRespMsg(-1, "complete upload failed", nil).JSONBytes())
		return
	}
	totalCount := 0
	chunkCount := 0
	for i := 0; i < len(data); i++ {
		k := string(data[i].([]byte))
		v := string(data[i+1].([]byte))
		if k == "chunkcount" {
			totalCount, _ = strconv.Atoi(v)
		} else if strings.HasPrefix(k, "chkidx_") && v == "1" {
			chunkCount++
		}
	}
	if totalCount != chunkCount {
		w.Write(util.NewRespMsg(-2, "invalid request", nil).JSONBytes())
		return
	}
	//4 TODO 合并分块
	//5 更新唯一文件表及用户文件表
	fsize, _ := strconv.Atoi(filesize)
	dblayer.OnFileUploadFinished(filehash, filename, int64(fsize), "")
	dblayer.OnUserFileUploadFinished(username, filehash, filename, int64(fsize))
	//6 返回处理结果到客户端
	w.Write(util.NewRespMsg(0, "OK", nil).JSONBytes())
}

//取消分块上传
func CancelUploadPartHandler(w http.ResponseWriter, r *http.Request){
	//1 删除已存在的分块状态
	//2 删除redis缓存状态
	//3 更新mysql文件status
}

func MultipartUploadStatusHandler(w http.ResponseWriter, r *http.Request){
	//1 检查分块上传状态是否有效
	//2 获取分块初始化信息
	//3 获取已上传分块信息
}