package rabbitmq

import (
	"fmt"
	"github.com/streadway/amqp"
)

func Publish(ex, routingKey string, msg []byte) bool {
	if !initConnection() {
		return false
	}
	err := mq.Channel.Publish(ex, routingKey, false, false, amqp.Publishing{ContentType: "text/plain", Body: msg})
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}