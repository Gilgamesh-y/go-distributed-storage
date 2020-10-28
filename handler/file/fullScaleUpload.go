package file

import (
	"DistributedStorage/conf"
	"DistributedStorage/fileMeta"
	"DistributedStorage/model/file_model"
	"DistributedStorage/mq/rabbitmq"
	"DistributedStorage/response"
	"encoding/json"
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
	storeType := viper.GetInt("store_type")
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
		if storeType == conf.StoreOSS {
			//file, _ := fileHeader.Open()
			ossPath := "full_scale/" + nowtime + "/" + fm.Hash + fileHeader.Filename
			//err =oss.Bucket().PutObject(ossPath, file)
			//if err != nil {
			//	fmt.Println(err)
			//	response.Resp(c, err, fm)
			//	return
			//}
			data := fileMeta.TransferData{
				Hash: fm.Hash,
				TmpPath: fm.Path,
				TargetPath: ossPath,
				StoreType: conf.StoreOSS,
			}
			publishData, _ := json.Marshal(data)
			success := rabbitmq.Publish(viper.GetString("TransExchangeName"), viper.GetString("TransOSSRoutingKey"), publishData)
			if !success {
				// TODO: retry
			}
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
