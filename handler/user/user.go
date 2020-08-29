package user

import (
	"DistributedStorage/model/user_model"
	"DistributedStorage/response"
	"github.com/gin-gonic/gin"
	"strconv"
)

func AddUser(c *gin.Context) {
	var user user_model.User
	if err := c.ShouldBind(&user); err != nil {
		response.Resp(c, err, user)
		return
	}
	err := user_model.AddUser(&user)

	response.Resp(c, err, user)
}

func GetUser(c *gin.Context) {
	var user user_model.User
	id, _ := strconv.Atoi(c.Param("id"))
	user.Id = int64(id)
	err := user_model.GetUser(&user)

	response.Resp(c, err, user)
}