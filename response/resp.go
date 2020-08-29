package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ReturnData struct {
	Code	int `json:"code"`
	Message string `json:"message"`
	Data	interface{} `json:"data"`
}

func Resp(c *gin.Context, err error, data interface{}) {
	code, message := formatErr(err)
	c.JSON(http.StatusOK, ReturnData{
		Code:	code,
		Message:message,
		Data:	data,
	})

	return
}