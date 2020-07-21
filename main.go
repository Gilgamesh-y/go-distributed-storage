package main

import (
	"DistributedStorage/conf"
	"DistributedStorage/route"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var cfg = pflag.StringP("conf", "c", "", "go-distributed-strorage config file path")

func main() {
	pflag.Parse()
	if err := conf.Init(*cfg); err != nil {
		panic(err)
	}

	gin.SetMode("debug")
	g := gin.New()
	route.Load(g)
	g.Run(viper.GetString("port"))
}