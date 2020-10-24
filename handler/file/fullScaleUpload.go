package file

import (
	"DistributedStorage/fileMeta"
	"DistributedStorage/model/file_model"
	"DistributedStorage/response"
	"DistributedStorage/store/oss"
	"fmt"
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
	nowtime := time.Now().Format("2006-01-02")
	uploadDir := pwd + viper.GetString("upload_dir") + nowtime + "/upload/"
	for _, fileHeader := range files {
		fm := &fileMeta.FileMeta{
			Name: fileHeader.Filename,
			Size: fileHeader.Size,
			UpdatedAt: nowtime,
		}
		fm.FileNameToSha1()
		fm.Path = uploadDir + fm.Hash
		err := fm.CreateDirIfNotExist(uploadDir)
		if err != nil {
			response.Resp(c, err, fm)
			return
		}

		// Save to local storage
		err = c.SaveUploadedFile(fileHeader, fm.Path)
		if err != nil {
			response.Resp(c, err, fm)
			return
		}


		// Save to ali oss
		file, _ := fileHeader.Open()
		fm.Path = "full_scale/" + nowtime + "/" + fm.Hash + fileHeader.Filename
		err =oss.Bucket().PutObject(fm.Path, file)
		if err != nil {
			fmt.Println(err)
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
