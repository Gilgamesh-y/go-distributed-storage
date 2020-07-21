package route

import "github.com/gin-gonic/gin"

func Options(c *gin.Context)  {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
	c.Header("Allow", "HEAD, GET, POST, PUT, PATCH, DELETE, OPTIONS")
	c.Header("Content-Type", "application/json")
	c.Next()
}