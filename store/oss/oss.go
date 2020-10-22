package oss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/viper"
)

var client *oss.Client

func Client() *oss.Client {
	if client == nil {
		client, err := oss.New(viper.GetString("endPoint"),
			viper.GetString("accessKeyId"),
			viper.GetString("accessKeySecret"))
		if err != nil {
			fmt.Println(err)
		}
		// 获取存储空间。
		return client
	}
	return client
}

func Bucket() (*oss.Bucket, error) {
	// 获取存储空间。
	return Client().Bucket(viper.GetString("bucketName"))
}

func GetOssFileList(nextMarker string, page int) (oss.ListObjectsResult, error) {
	bucket, _ := Bucket()
	marker := oss.Marker(nextMarker)
	list, err := bucket.ListObjects(oss.MaxKeys(page), marker)
	if err != nil {
		return list, err
	}

	return list, nil
}

func DeleteOssFile(name string) (int, error) {
	bucket, _ := Bucket()
	err := bucket.DeleteObject(name)
	if err != nil {
		return 0, err
	}

	return 1, nil
}

func CountOssFile() (int, error) {
	bucket, _ := Bucket()
	var lsRes oss.ListObjectsResult
	var err error

	lsRes, err = bucket.ListObjects(oss.MaxKeys(1000))
	if err != nil {
		return 0, err
	}
	count := len(lsRes.Objects)

	return count, nil
}