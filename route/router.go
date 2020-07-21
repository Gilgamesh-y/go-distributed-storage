package route

import (
	"DistributedStorage/handler"
	"github.com/gin-gonic/gin"
)

func Load(g *gin.Engine) *gin.Engine {
	g.Use(gin.Recovery())
	g.Use(Options)

	g.POST("upload", handler.Upload)
	g.GET("files", handler.Upload)

	return g
}