package route

import (
	"DistributedStorage/handler/file"
	"DistributedStorage/handler/user"
	"github.com/gin-gonic/gin"
)

func Load(g *gin.Engine) *gin.Engine {
	g.Use(gin.Recovery())
	g.Use(Options)

	g.POST("upload", file.Upload)
	g.GET("files", file.Get)

	g.POST("add_user", user.AddUser)
	g.GET("get_user/:id", user.GetUser)

	return g
}