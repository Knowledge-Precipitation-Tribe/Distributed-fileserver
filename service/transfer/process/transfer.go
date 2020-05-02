package process

import (
	"Distributed-fileserver/service/transfer/customLog"
	"bufio"
	"encoding/json"
	"go.uber.org/zap"
	"log"
	"os"

	"Distributed-fileserver/mq"
	dbcli "Distributed-fileserver/service/dbproxy/client"
	"Distributed-fileserver/store/oss"
)

// Transfer : 处理文件转移
func Transfer(msg []byte) bool {
	log.Println(string(msg))

	pubData := mq.TransferData{}
	err := json.Unmarshal(msg, &pubData)
	if err != nil {
		customLog.Logger.Error("处理文件转移json解析失败", zap.Error(err))
		return false
	}

	fin, err := os.Open(pubData.CurLocation)
	if err != nil {
		customLog.Logger.Error("处理文件转移os open失败", zap.Error(err))
		return false
	}

	err = oss.Bucket().PutObject(
		pubData.DestLocation,
		bufio.NewReader(fin))
	if err != nil {
		customLog.Logger.Error("处理文件转移oss.bucket.put失败", zap.Error(err))
		return false
	}

	resp, err := dbcli.UpdateFileLocation(
		pubData.FileHash,
		pubData.DestLocation)
	if err != nil {
		customLog.Logger.Error("处理文件转移uploadFileLocation失败", zap.Error(err))
		return false
	}
	if !resp.Suc {
		customLog.Logger.Error("处理文件转移失败", zap.String("更新数据库异常，请检查:",pubData.FileHash))
		return false
	}
	return true
}
