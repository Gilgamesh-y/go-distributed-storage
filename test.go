package main

import "github.com/gin-gonic/gin"

func main() {
	gin.SetMode("debug")
	r := gin.New()
	r.GET("ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	r.Run(":8080")
}