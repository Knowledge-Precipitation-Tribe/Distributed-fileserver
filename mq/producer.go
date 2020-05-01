package mq

import (
	"Distributed-fileserver/config"
	log "Distributed-fileserver/zaplogger"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

//发布消息
func Publish(exchange string, routingKey string, msg []byte) bool {
	//检查channel是否正常
	if !initChannel(config.RabbitURL) {
		return false
	}
	//执行消息发布
	err := channel.Publish(
		exchange,
		routingKey,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        msg,
		})
	if err != nil {
		log.GetLogger().Error(
			"error mq Publish",
			zap.Error(err))
		return false
	}
	return true
}
