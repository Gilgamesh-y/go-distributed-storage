package handler

import (
	"DistributedStorage/fileMeta"
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
		filemeta := fileMeta.FileMeta{
			Name: file.Filename,
			Path: pwd + viper.GetString("upload_dir") + file.Filename,
			UpdateAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		err := c.SaveUploadedFile(file, filemeta.Path)
		if err != nil {
			panic(err)
		}
		filemeta.Size = filemeta.GetSize()
	}
}
