package route

import (
	"DistributedStorage/handler/file"
	"DistributedStorage/handler/user"
	"github.com/gin-gonic/gin"
)

func Load(g *gin.Engine) *gin.Engine {
	g.Use(gin.Recovery())
	g.Use(Options)

	g.GET("files", file.Get)

	g.POST("add_user", user.AddUser)
	g.GET("get_user/:id", user.GetUser)

	f := g.Group("/file")
	{
		fs := f.Group("/single")
		{
			fs.POST("upload", file.Upload)
		}
		fm := f.Group("/multipart")
		{
			fm.POST("init", file.InitMultipartUploadInfo)
			fm.POST("upload", file.MultipartUpload)
			fm.POST("upload_complete", file.MultipartUploadComplete)
		}
	}

	return g
}