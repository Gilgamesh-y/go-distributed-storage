package rabbitmq

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	Conn *amqp.Connection
	Channel *amqp.Channel
	NotifySwitch chan *amqp.Error
}

var mq *RabbitMQ

func init() {
	// 是否开启异步转移功能，开启时才初始化rabbitMQ连接
	if !viper.GetBool("rabbitmqAsyncSend") {
		return
	}
	if initConnection() {
		mq.Channel.NotifyClose(mq.NotifySwitch)
	}
	// 断线自动重连
	go func() {
		for {
			select {
			case msg := <-mq.NotifySwitch:
				mq.Conn = nil
				mq.Channel = nil
				fmt.Println("channel closed: %+v\n", msg)
				initConnection()
			}
		}
	}()
}

func initConnection() bool {
	if mq.Conn != nil {
		return true
	}

	conn, err := amqp.Dial(viper.GetString("rabbitmqUrl"))
	if err != nil {
		fmt.Println(err)
		return false
	}
	mq.Conn = conn

	channel, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
		return false
	}
	mq.Channel = channel

	return true
}