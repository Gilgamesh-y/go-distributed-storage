package main

// 测试分块上传 (场景：正常完成上传)

import (
	"DistributedStorage/fileMeta"
	"bufio"
	"encoding/json"
	"io"
	"path/filepath"

	//"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	//"path/filepath"
	"strconv"

	jsonit "github.com/json-iterator/go"
)

const (
	apiHost           = "http://localhost:8080/"
	apiUploadInit     = apiHost + "file/multipart/init"
	apiUploadPart     = apiHost + "file/multipart/upload"
	apiUploadComplete = apiHost + "file/mpupload/complete"
	apiUploadCancel   = apiHost + "file/mpupload/cancel"
)

// 当前测试只上传的分片数量, 默认为全部上传
var uploadChunkCount int

// MultipartUploadInfo : 初始化的分片信息
type MultipartUploadInfo struct {
	FileHash   string
	FileSize   int
	//UploadId   string
	ChunkSize  int
	ChunkCount int
	// 已经存在的分块，告诉客户端可以跳过这些分块，无需重复上传
	ChunkExists []int
}

// UploadInitResponse : 初始化接口返回的数据
type UploadInitResponse struct {
	Code int                 `json:code`
	Msg  string              `json:msg`
	Data MultipartUploadInfo `json:data`
}

// 实际上传分块逻辑
// chunkIdxs : 实际需要上传的分块
func uploadPartsSpecified(filename string, targetURL string, chunkSize int, chunkIdxs []int, uploadId string) error {
	f, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()

	bfRd := bufio.NewReader(f)
	index := 0

	ch := make(chan int)
	buf := make([]byte, chunkSize) //每次读取chunkSize大小的内容
	for {
		n, err := bfRd.Read(buf)
		if n <= 0 {
			break
		}
		index++

		// 判断当前所在的块是否需要上传
		if contained, err := fileMeta.Contain(chunkIdxs, index); err != nil || !contained {
			continue
		}

		// 可以不使用bufCopied, 直接传将buf slice作为参数传递(值传递)进去计算sha1
		bufCopied := make([]byte, 5*1048576)
		copy(bufCopied, buf)

		go func(b []byte, curIdx int) {
			hash := fileMeta.Sha1(b)
			resp, err := http.PostForm(
				targetURL,
				url.Values{
					"chunk_index": {strconv.Itoa(curIdx)},
					"upload_id": {uploadId},
					"hash": {hash},
				})
			if err != nil {
				fmt.Println(err)
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("%+v %+v\n", string(body), err)
			}
			resp.Body.Close()

			ch <- curIdx
		}(bufCopied[:n], index)

		//遇到任何错误立即返回，并忽略 EOF 错误信息
		if err != nil {
			if err == io.EOF {
				break
			} else {
				fmt.Println(err.Error())
			}
		}
	}

	for idx := 0; idx < len(chunkIdxs); idx++ {
		select {
		case res := <-ch:
			fmt.Printf("完成传输块index: %d\n", res)
		}
	}

	fmt.Printf("全部完成以下分块传输: %+v\n", chunkIdxs)
	return nil
}

func main() {
	uploadFilePath := "/Users/wrath/Documents/myspace/go-distributed-storage/storage/apache-tomcat-8.5.58.zip"
	fm := &fileMeta.FileMeta{
		Path: uploadFilePath,
	}
	// 判断参数是否有效
	if exist := fm.IsExist(); !exist {
		fmt.Println("Error: 无效文件路径，请检查")
		return
	}

	filesize := fm.GetSize()
	fm.FileNameToSha1()
	// 1. 请求初始化分块上传接口
	resp, err := http.PostForm(
		apiUploadInit,
		url.Values{
			"hash": {fm.Hash},
			"file_size": {strconv.FormatInt(filesize, 10)},
		})

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	// 2. 得到uploadID以及服务端指定的分块大小chunkSize
	uploadID := jsonit.Get(body, "data").Get("UploadId").ToString()
	chunkSize := jsonit.Get(body, "data").Get("ChunkSize").ToInt()
	fmt.Printf("uploadid: %s  chunksize: %d\n", uploadID, chunkSize)

	// 3. 请求分块上传接口
	var initResp UploadInitResponse
	err = json.Unmarshal(body, &initResp)
	if err != nil {
		fmt.Printf("Parse error: %s\n", err.Error())
		os.Exit(-1)
	}
	var chunksToUpload []int
	for idx := 1; idx <= initResp.Data.ChunkCount; idx++ {
		chunksToUpload = append(chunksToUpload, idx)
	}
	uploadChunkCount = len(chunksToUpload)
	tURL := apiUploadPart
	// 上传所有分块
	uploadPartsSpecified(uploadFilePath, tURL, chunkSize, chunksToUpload, uploadID)

	// 4. 请求分块完成接口
	resp, err = http.PostForm(
		apiUploadComplete,
		url.Values{
			"filehash": {fhash},
			"filesize": {strconv.Itoa(filesize)},
			"filename": {filepath.Base(uploadFilePath)},
			"uploadid": {uploadID},
		})

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	defer resp.Body.Close()
	// 5. 打印分块上传结果
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	fmt.Printf("complete result: %s\n", string(body))
}