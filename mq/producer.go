package mq

import (
	"Distributed-fileserver/config"
	log "Distributed-fileserver/zaplogger"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

var conn *amqp.Connection
var channel *amqp.Channel

//初始化channel
func initChannel() bool{
	if channel != nil{
		return true
	}
	//获得rabbitmq的连接
	conn, err := amqp.Dial(config.RabbitURL)
	if err != nil{
		log.GetLogger().Error(
			"error initChannel get rabbitmq conn",
			zap.Error(err))
		return false
	}
	//打开一个channel，用于消息发布与接收
	channel, err = conn.Channel()
	if err != nil{
		log.GetLogger().Error(
			"error initChannel get channel",
			zap.Error(err))
		return false
	}
	return true
}

