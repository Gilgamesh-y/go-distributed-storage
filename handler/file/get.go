package file

import (
	"DistributedStorage/model/file_model"
	"DistributedStorage/response"
	"github.com/gin-gonic/gin"
)

func Get(c *gin.Context) {
	fm, err := file_model.Get()
	if err != nil {
		response.Resp(c, err, fm)
		return
	}
	response.Resp(c, err, fm)
}