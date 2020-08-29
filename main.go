package main

import (
	"DistributedStorage/cache"
	"DistributedStorage/conf"
	"DistributedStorage/model"
	"DistributedStorage/route"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var cfg = pflag.StringP("conf", "c", "", "go-distributed-strorage config file_model path")

func main() {
	pflag.Parse()
	if err := conf.Init(*cfg); err != nil {
		panic(err)
	}

	model.Init()

	gin.SetMode("debug")
	g := gin.New()
	route.Load(g)

	cache.Init()

	err := cache.Set("SET", "name", "wrath")
	m, err := cache.GetString("name")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(m)
	g.Run(":"+viper.GetString("port"))
}