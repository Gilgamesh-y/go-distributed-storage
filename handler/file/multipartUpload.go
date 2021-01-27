package file

import (
	"DistributedStorage/cache"
	"DistributedStorage/fileMeta"
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
	ChunkSize  int64  `form:"chunk_size"`
	ChunkCount int    `form:"chunk_count"`
	chunkIndexExists []int `form:"chunk_index_exists"`
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

const (
	RedisPrefixKey = "mpu_"
)

/**
 * Init the information about multipart upload
 */
func InitMultipartUploadInfo(c *gin.Context) {
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

	// Does the file exist
	fileExists, _ := redis.Bool(cache.Get("EXISTS", RedisPrefixKey+impu.Hash))
	if fileExists {
		uploadId, _ := redis.String(cache.Get("GET", RedisPrefixKey+impu.Hash))
		impu.UploadId, _ = strconv.ParseInt(uploadId, 10, 64)
	}

	/**
	 * If uploading for the first time, create a new upload_id, if not, get the uploaded parts according to upload_id
	 */
	chunkIndexExists := []int{}
	if impu.UploadId == 0 {
		impu.UploadId = worker.GetId()
	} else {
		chunks, _ := redis.Values(cache.Get("HGETALL", RedisPrefixKey+strconv.FormatInt(impu.UploadId, 10)))
		for i := 0; i < len(chunks); i += 2 {
			k := string(chunks[i].([]byte))
			v := string(chunks[i+1].([]byte))
			if strings.HasPrefix(k, "chunk_index") && v == "1" {
				// Get the value about chunk_index from key
				chunkIndex, _ := strconv.Atoi(k[11:len(k)])
				chunkIndexExists = append(chunkIndexExists, chunkIndex)
			}
		}
		impu.chunkIndexExists = chunkIndexExists
	}

	// Get the information about multipart upload
	impu.ChunkSize = 5*1024*1024 // 5MB
	impu.ChunkCount = int(impu.FileSize/impu.ChunkSize)

	if len(chunkIndexExists) == 0 {
		// Save the information of the file into redis
		key := RedisPrefixKey+strconv.FormatInt(impu.UploadId, 10)
		cache.Set("HSET", key, "chunk_count", impu.ChunkCount, "EX", 7 * 86400)
		cache.Set("HSET", key, "hash", impu.Hash, "EX", 7 * 86400)
		cache.Set("HSET", key, "file_size", impu.FileSize, "EX", 7 * 86400)
		// Save the upload_id of the file hash
		cache.Set("SET", RedisPrefixKey+impu.Hash, impu.UploadId, "EX", 7 * 86400)
	}

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
	nowtime := time.Now().Format("2006-01-02")
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

	key := RedisPrefixKey+strconv.FormatInt(mpu.UploadId, 10)
	cache.Set("HSET", key, "chunk_index_" + strconv.Itoa(mpu.ChunkIndex), 1, "EX", 7 * 86400)
	response.Resp(c, nil, mpu)
}

/**
 * Notice to upload and merge
 */
func MultipartUploadComplete(c *gin.Context) {
	var mpuc MultipartUploadCompleteStruct
	if err := c.ShouldBind(&mpuc); err != nil {
		response.Resp(c, err, mpuc)
		return
	}
	// Determine whether all the chunks are uploaded
	mpuData, err := redis.Values(cache.Get("HGETALL", RedisPrefixKey + strconv.FormatInt(mpuc.UploadId, 10)))
	if err != nil {
		response.Resp(c, err, "上传失败")
		return
	}
	totalCount := 0
	chunkCount := 0
	// HGETALL's key and value is in  array for why i += 2
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

/**
 * Notice to cancel upload
 */
func CancelUpload(c *gin.Context) {
	hash := c.PostForm("hash")
	// TODO delete redis cache
	uploadId, _ := redis.String(cache.Get("GET", RedisPrefixKey + hash))
	if uploadId == "" {
		response.Resp(c, response.New(response.UploadIdNotFound, hash),  nil)
		return
	}
	// TODO delete existing chunked files
	pwd, _ := os.Getwd()
	nowtime := time.Now().Format("2006-01-02")
	uploadDir := pwd + viper.GetString("upload_dir") + nowtime + "/multipart_upload/" + uploadId
	if !fileMeta.DelFileByShell(uploadDir) {
		response.Resp(c, response.New(response.DelFileFail, uploadId),  nil)
		return
	}
	// TODO update mysql
}

/**
 * Get the information of the upload status
 */
func MultipartUploadStatus(c *gin.Context) {
	// TODO get unsuccessful data from redis according to upload_id
}