package mq

import (
	log "Distributed-fileserver/zaplogger"
	"go.uber.org/zap"
)

var done chan bool

//开始监听队列，获取消息
func StratConsume(qName string, cName string, callback func(msg []byte) bool) {
	msgs, err := channel.Consume(
		qName,
		cName,
		true,  //是否自动恢复ack信号给消息发布者
		false, //制定当前是否为唯一的消费者
		false, //没用
		false, //是否等待
		nil)
	if err != nil {
		log.GetLogger().Error(
			"error initChannel get rabbitmq conn",
			zap.Error(err))
		return
	}

	done = make(chan bool)

	go func() {
		//循环获取队列的消息
		for msg := range msgs {
			//调用callback获取返回的消息
			result := callback(msg.Body)
			if !result {
				//Todo 将失败的人物写到另一个队列，用于异常情况的消息重试
			}
		}
	}()

	//没有消息的话会一直阻塞
	<-done

	//关闭rabbitMQ的channel
	channel.Close()
}

// StopConsume : 停止监听队列
func StopConsume() {
	done <- true
}
