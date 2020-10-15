package file

import (
	"DistributedStorage/cache"
	"DistributedStorage/fileMeta"
	"DistributedStorage/model/file_model"
	"DistributedStorage/response"
	"DistributedStorage/util/snowflake"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"strings"
	"time"
)

type InitMultiPartUploadStruct struct {
	UploadId   int64  `form:"upload_id"`
	Hash       string `form:"hash" binding:"required"`
	FileSize   int64  `form:"file_size" binding:"required"`
	ChunkSize  int64  `form:"chunk_size" binding:"required"`
	ChunkCount int    `form:"chunk_count" binding:"required"`
}

type MultiPartUploadStruct struct {
	UploadId   int64 `form:"upload_id" binding:"required"`
	ChunkIndex int   `form:"chunk_index" binding:"required"`
}

type MultipartUploadCompleteStruct struct {
	UploadId int64 `form:"upload_id" binding:"required"`
	Hash     int   `form:"hash" binding:"required"`
	FileSize int64 `form:"file_size" binding:"required"`
	FileName int   `form:"file_name" binding:"required"`
}

/**
 * Upload a file
 */
func Upload(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["files"]
	pwd, _ := os.Getwd()
	nowtime := time.Now().Format("2006-01-02 15:04:05")
	uploadDir := pwd + viper.GetString("upload_dir") + nowtime + "/upload/"
	for _, file := range files {
		fm := &fileMeta.FileMeta{
			Name: file.Filename,
			Path: uploadDir + file.Filename,
			Size: file.Size,
			UpdatedAt: nowtime,
		}
		fm.ToSha1()
		err := fm.CreateDirIfNotExist(uploadDir)
		if err != nil {
			response.Resp(c, err, fm)
			return
		}
		err = c.SaveUploadedFile(file, fm.Path)
		if err != nil {
			response.Resp(c, err, fm)
			return
		}

		existFm, err := file_model.GetByHash(fm.Hash)
		if err != nil {
			response.Resp(c, err, fm)
			return
		}
		if existFm.Id != 0 {
			response.Resp(c, response.FileExist, existFm)
			return
		}

		_, err = file_model.Insert(fm)
		if err != nil {
			response.Resp(c, err, fm)
			return
		}
	}
	response.Resp(c, nil, nil)
}

/**
 * Init the information about multipart upload
 */
func InitMultiPartUploadInfo(c *gin.Context) {
	var impu InitMultiPartUploadStruct
	if err := c.ShouldBind(&impu); err != nil {
		response.Resp(c, err, impu)
		return
	}
	worker, err := snowflake.NewWorker(1)
	if err != nil {
		response.Resp(c, err, nil)
		return
	}
	// Get the information about multipart upload
	impu.UploadId = worker.GetId()
	impu.ChunkSize = 5*1024*1024 // 5MB
	impu.ChunkCount = int(impu.FileSize/impu.ChunkSize)

	// Save the information of the file into redis
	key := "mpu_"+strconv.FormatInt(impu.UploadId, 10)
	cache.Set("HSET", key, "chunk_count", impu.ChunkCount)
	cache.Set("HSET", key, "hash", impu.Hash)
	cache.Set("HSET", key, "file_size", impu.FileSize)

	response.Resp(c, nil, impu)
}

/**
 * Save the part of the file
 */
func MultipartUpload(c *gin.Context) {
	var mpu MultiPartUploadStruct
	if err := c.ShouldBind(&mpu); err != nil {
		response.Resp(c, err, mpu)
		return
	}
	pwd, _ := os.Getwd()
	nowtime := time.Now().Format("2006-01-02 15:04:05")
	uploadDir := pwd + viper.GetString("upload_dir") + nowtime + "/multipart_upload/" + strconv.FormatInt(mpu.UploadId, 10) + "/"
	fm := &fileMeta.FileMeta{
		Path: uploadDir + strconv.Itoa(mpu.ChunkIndex),
	}

	// TODO Verify the hash value

	// Save the content of the chunk
	err := fm.CreateDirIfNotExist(uploadDir)
	if err != nil {
		response.Resp(c, err, fm)
		return
	}
	fd, err := os.Create(fm.Path)
	if err != nil {
		response.Resp(c, err, fm)
		return
	}
	defer fd.Close()
	buf := make([]byte, 1024*1024)
	for {
		n, err := c.Request.Body.Read(buf)
		fd.Write(buf[:n])
		if err != nil {
			break
		}
	}

	key := "mpu_"+strconv.FormatInt(mpu.UploadId, 10)
	cache.Set("HSET", key, "chunk_index_" + strconv.Itoa(mpu.ChunkIndex), 1)
	response.Resp(c, nil, nil)
}

func MultipartUploadComplete(c *gin.Context) {
	var mpuc MultipartUploadCompleteStruct
	if err := c.ShouldBind(&mpuc); err != nil {
		response.Resp(c, err, mpuc)
		return
	}
	// Determine whether all the chunks are uploaded
	mpuData, err := redis.Values(cache.Get("HGETALL", "mpu_" + strconv.FormatInt(mpuc.UploadId, 10)))
	if err != nil {
		response.Resp(c, err, "上传失败")
		return
	}
	totalCount := 0
	chunkCount := 0
	for i := 0; i < len(mpuData); i += 2 {
		key := string(mpuData[i].([]byte))
		val := string(mpuData[i + 1].([]byte))
		if key == "chunk_count" {
			totalCount, _ = strconv.Atoi(val)
		}
		if strings.HasPrefix(key, "chunk_index") && val == "1" {
			chunkCount += 1
		}
		if totalCount != chunkCount {
			response.Resp(c, err, "上传失败")
			return
		}
	}

	// TODO Merge chunk

	// TODO Update database

}