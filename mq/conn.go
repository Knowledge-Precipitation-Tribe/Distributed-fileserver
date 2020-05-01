package mq

import (
	"Distributed-fileserver/config"
	log "Distributed-fileserver/zaplogger"
	"fmt"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

var conn *amqp.Connection
var channel *amqp.Channel

// 如果异常关闭，会接收通知
var notifyClose chan *amqp.Error

// UpdateRabbitHost : 更新mq host
func UpdateRabbitHost(host string) {
	config.RabbitURL = host
}

// Init : 初始化MQ连接信息
func Init() {
	// 是否开启异步转移功能，开启时才初始化rabbitMQ连接
	if !config.AsyncTransferEnable {
		return
	}
	if initChannel(config.RabbitURL) {
		channel.NotifyClose(notifyClose)
	}
	// 断线自动重连
	go func() {
		for {
			select {
			case msg := <-notifyClose:
				conn = nil
				channel = nil
				log.GetLogger().Info(fmt.Sprintf("onNotifyChannelClosed: %+v\n", msg))
				initChannel(config.RabbitURL)
			}
		}
	}()
}

//初始化channel
func initChannel(rabbitHost string) bool {
	if channel != nil {
		return true
	}

	var err error
	conn, err = amqp.Dial(rabbitHost)
	if err != nil {
		log.GetLogger().Error(
			"error initChannel get rabbitmq conn",
			zap.Error(err))
		return false
	}

	channel, err = conn.Channel()
	if err != nil {
		log.GetLogger().Error(
			"error initChannel get channel",
			zap.Error(err))
		return false
	}
	return true
}
