package oss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/spf13/viper"
)

var client *oss.Client

func Client() *oss.Client {
	if client == nil {
		client, err := oss.New(viper.GetString("end_point"),
			viper.GetString("access_key_id"),
			viper.GetString("access_key_secret"))
		if err != nil {
			fmt.Println(err)
		}
		// 获取存储空间。
		return client
	}
	return client
}

func Bucket() *oss.Bucket {
	// 获取存储空间。
	bucket, err := Client().Bucket(viper.GetString("bucket_name"))
	if err != nil {
		fmt.Println(err)
	}
	return bucket
}

func GetOssFileList(nextMarker string, page int) (oss.ListObjectsResult, error) {
	bucket := Bucket()
	marker := oss.Marker(nextMarker)
	list, err := bucket.ListObjects(oss.MaxKeys(page), marker)
	if err != nil {
		return list, err
	}

	return list, nil
}

func DeleteOssFile(name string) (int, error) {
	bucket := Bucket()
	err := bucket.DeleteObject(name)
	if err != nil {
		return 0, err
	}

	return 1, nil
}

func CountOssFile() (int, error) {
	bucket := Bucket()
	var lsRes oss.ListObjectsResult
	var err error

	lsRes, err = bucket.ListObjects(oss.MaxKeys(1000))
	if err != nil {
		return 0, err
	}
	count := len(lsRes.Objects)

	return count, nil
}