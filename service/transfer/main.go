package main

import (
	"DistributedStorage/fileMeta"
	"DistributedStorage/mq/rabbitmq"
	"DistributedStorage/store/oss"
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"os"
)

func ProcessTransfer(msg []byte) bool {
	publishData := fileMeta.TransferData{}
	err := json.Unmarshal(msg, publishData)
	if err != nil {
		fmt.Println(err)
		return false
	}
	file, err := os.Open(publishData.TmpPath)
	if err != nil {
		fmt.Println(err)
		return false
	}

	err = oss.Bucket().PutObject(publishData.TargetPath, bufio.NewReader(file))
	if err != nil {
		fmt.Println(err)
		return false
	}

	// TODO: update mysql

	return true
}

func main() {
	rabbitmq.StartConsume(viper.GetString("TransOSSQueueName"), "transfer_oss", ProcessTransfer)
}