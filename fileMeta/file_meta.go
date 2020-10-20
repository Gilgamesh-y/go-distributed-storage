package fileMeta

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type FileMeta struct {
	Id int
	Name string
	Size int64
	UpdatedAt string
	Path string
	Hash string
	UploadId string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

func Set(fm FileMeta) {
	fileMetas[fm.Name] = fm
}

func (fm FileMeta) Get() FileMeta {
	return fileMetas[fm.Name]
}

func (fm FileMeta) GetSize() int64 {
	file, _ := os.Stat(fm.Path)
	return file.Size()
}

func (fm FileMeta) GetName() string {
	file, _ := os.Stat(fm.Path)
	return file.Name()
}

func (fm FileMeta) IsExist() bool {
	_, err := os.Stat(fm.Path)
	if (os.IsNotExist(err)) {
		return false
	}
	return true
}

func (fm FileMeta) CreateDirIfNotExist(dir string) error {
	 if !fm.IsExist() {
		 err := os.MkdirAll(dir, 0777)
		 if err != nil {
			 return err
		 }
	 }
	 return nil
}

func (fm FileMeta) GetModTime() time.Time {
	file, _ := os.Stat(fm.Path)
	return file.ModTime()
}

func (fm *FileMeta) FileNameToSha1() {
	if fm.Name == "" {
		fm.Name = fm.GetName()
	}
	s := fm.Name + strconv.FormatInt(fm.Size, 10)
	h := sha1.New()
	h.Write([]byte(s))
	fm.Hash = hex.EncodeToString(h.Sum([]byte("")))
}

// Contain : 判断某个元素是否在 slice,array ,map中
func Contain(target interface{}, obj interface{}) (bool, error) {
	targetVal := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		// 是否在slice/array中
		for i := 0; i < targetVal.Len(); i++ {
			if targetVal.Index(i).Interface() == obj {
				return true, nil
			}
		}
	case reflect.Map:
		// 是否在map key中
		if targetVal.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true, nil
		}
	default:
		fmt.Println(reflect.TypeOf(target).Kind())
	}

	return false, errors.New("not in this array/slice/map")
}

func Sha1(data []byte) string {
	_sha1 := sha1.New()
	_sha1.Write(data)
	return hex.EncodeToString(_sha1.Sum([]byte("")))
}

// RemovePathByShell : 通过调用shell来删除制定目录 合并成功将返回true, 否则返回false
func DelFileByShell(path string) bool {
	cmdStr := strings.Replace(`
	#!/bin/bash
	chunkDir="/data/chunks/"
	targetDir=$1
	# 增加条件判断，避免误删  (指定的路径包含且不等于chunkDir)
	if [[ $targetDir =~ $chunkDir ]] && [[ $targetDir != $chunkDir ]]; then 
	  rm -rf $targetDir
	fi
	`, "$1", path, 1)
	delCmd := exec.Command("bash", "-c", cmdStr)
	if _, err := delCmd.Output(); err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

