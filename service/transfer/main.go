package main

import (
	"Distributed-fileserver/service/transfer/customLog"
	"go.uber.org/zap"
	"log"
	"time"

	"Distributed-fileserver/common"
	"Distributed-fileserver/config"
	"Distributed-fileserver/mq"
	dbproxy "Distributed-fileserver/service/dbproxy/client"
	"Distributed-fileserver/service/transfer/process"

	"github.com/micro/cli"
	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/registry/consul"
	_ "github.com/micro/go-plugins/registry/kubernetes"
)

func startRPCService() {
	service := micro.NewService(
		micro.Name("go.micro.service.transfer"), // 服务名称
		micro.RegisterTTL(time.Second*10),       // TTL指定从上一次心跳间隔起，超过这个时间服务会被服务发现移除
		micro.RegisterInterval(time.Second*5),   // 让服务在指定时间内重新注册，保持TTL获取的注册时间有效
		micro.Flags(common.CustomFlags...),
	)
	service.Init(
		micro.Action(func(c *cli.Context) {
			// 检查是否指定mqhost
			mqhost := c.String("mqhost")
			if len(mqhost) > 0 {
				log.Println("custom mq address: " + mqhost)
				mq.UpdateRabbitHost(mqhost)
			}
		}),
	)

	// 初始化dbproxy client
	dbproxy.Init(service)

	if err := service.Run(); err != nil {
		customLog.Logger.Error("transfer main service run 失败", zap.Error(err))
	}
}

func startTranserService() {
	if !config.AsyncTransferEnable {
		log.Println("异步转移文件功能目前被禁用，请检查相关配置")
		return
	}
	log.Println("文件转移服务启动中，开始监听转移任务队列...")

	// 初始化mq client
	mq.Init()

	mq.StartConsume(
		config.TransOSSQueueName,
		"transfer_oss",
		process.Transfer)
}

func main() {
	// 文件转移服务
	go startTranserService()

	// rpc 服务
	startRPCService()
}
