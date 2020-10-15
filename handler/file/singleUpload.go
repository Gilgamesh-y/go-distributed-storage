package file

import (
	"DistributedStorage/fileMeta"
	"DistributedStorage/model/file_model"
	"DistributedStorage/response"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"os"
	"time"
)

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
