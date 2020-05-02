package handler

import (
	"Distributed-fileserver/service/account/customLog"
	"context"
	"encoding/json"
	"go.uber.org/zap"

	"Distributed-fileserver/common"
	proto "Distributed-fileserver/service/account/proto"
	dbcli "Distributed-fileserver/service/dbproxy/client"
)

// UserFiles : 获取用户文件列表
func (u *User) UserFiles(ctx context.Context, req *proto.ReqUserFile, res *proto.RespUserFile) error {
	dbResp, err := dbcli.QueryUserFileMetas(req.Username, int(req.Limit))
	if err != nil || !dbResp.Suc {
		customLog.Logger.Error("获取文件列表失败", zap.Error(err))
		res.Code = common.StatusServerError
		return err
	}

	userFiles := dbcli.ToTableUserFiles(dbResp.Data)
	data, err := json.Marshal(userFiles)
	if err != nil {
		customLog.Logger.Error("获取文件列表json解析失败", zap.Error(err))
		res.Code = common.StatusServerError
		return nil
	}

	res.FileData = data
	return nil
}

// UserFiles : 用户文件重命名
func (u *User) UserFileRename(ctx context.Context, req *proto.ReqUserFileRename, res *proto.RespUserFileRename) error {
	dbResp, err := dbcli.RenameFileName(req.Username, req.Filehash, req.NewFileName)
	if err != nil || !dbResp.Suc {
		customLog.Logger.Error("文件重命名失败", zap.Error(err))
		res.Code = common.StatusServerError
		return err
	}

	userFiles := dbcli.ToTableUserFiles(dbResp.Data)
	data, err := json.Marshal(userFiles)
	if err != nil {
		customLog.Logger.Error("文件重命名json解析失败", zap.Error(err))
		res.Code = common.StatusServerError
		return nil
	}

	res.FileData = data
	return nil
}
