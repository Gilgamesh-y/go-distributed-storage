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

func Upload(c *gin.Context) {
	form, _ := c.MultipartForm()
	files := form.File["files"]
	for _, file := range files {
		pwd, _ := os.Getwd()
		fm := fileMeta.FileMeta{
			Name: file.Filename,
			Path: pwd + viper.GetString("upload_dir") + file.Filename,
			UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		err := fm.CreateDirIfNotExist(pwd + viper.GetString("upload_dir"))
		if err != nil {
			response.Resp(c, err, fm)
			return
		}
		err = c.SaveUploadedFile(file, fm.Path)
		if err != nil {
			response.Resp(c, err, fm)
			return
		}
		fm.Size = fm.GetSize()
		_, err = file_model.Insert(fm)
		if err != nil {
			response.Resp(c, err, fm)
			return
		}
		response.Resp(c, nil, nil)
	}
}

func Get(c *gin.Context) {
	fm, err := file_model.Get()
	if err != nil {
		response.Resp(c, err, fm)
		return
	}
	response.Resp(c, err, fm)
}
