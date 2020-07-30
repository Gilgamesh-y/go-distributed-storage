package route

import (
	"DistributedStorage/handler/file"
	"github.com/gin-gonic/gin"
)

func Load(g *gin.Engine) *gin.Engine {
	g.Use(gin.Recovery())
	g.Use(Options)

	g.POST("upload", file.Upload)
	g.GET("files", file.Get)

	return g
}