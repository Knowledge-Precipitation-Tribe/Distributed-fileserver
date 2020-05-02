package main

import (
	"Distributed-fileserver/service/download/customLog"
	"go.uber.org/zap"
	"time"

	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/registry/consul"
	_ "github.com/micro/go-plugins/registry/kubernetes"

	"Distributed-fileserver/common"
	dbproxy "Distributed-fileserver/service/dbproxy/client"
	cfg "Distributed-fileserver/service/download/config"
	dlProto "Distributed-fileserver/service/download/proto"
	"Distributed-fileserver/service/download/route"
	dlRpc "Distributed-fileserver/service/download/rpc"
)

func startRPCService() {
	service := micro.NewService(
		micro.Name("go.micro.service.download"), // 在注册中心中的服务名称
		micro.RegisterTTL(time.Second*10),
		micro.RegisterInterval(time.Second*5),
		micro.Flags(common.CustomFlags...),
	)
	service.Init()

	// 初始化dbproxy client
	dbproxy.Init(service)

	dlProto.RegisterDownloadServiceHandler(service.Server(), new(dlRpc.Download))
	if err := service.Run(); err != nil {
		customLog.Logger.Error(" 从文件表查找记录失败", zap.Error(err))
	}
}

func startAPIService() {
	router := route.Router()
	router.Run(cfg.DownloadServiceHost)
}

func main() {
	// api 服务
	go startAPIService()

	// rpc 服务
	startRPCService()
}
