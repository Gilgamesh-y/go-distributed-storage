package rabbitmq

import "fmt"

var msgChan chan bool

func StartConsume(queueName, consumerName string, callback func(msg []byte) bool) {
	// Get message channel
	msgs, err := mq.Channel.Consume(queueName, consumerName, true, false, false, false, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	go func() {
		// Get message from queue
		for msg := range msgs {
			success := callback(msg.Body)
			if !success {
				// TODO: retry
			}
		}
	}()

	// If not exist new msg it will be block
	<- msgChan

	mq.Channel.Close()
}