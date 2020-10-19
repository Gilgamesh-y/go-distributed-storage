package main

import (
	"DistributedStorage/cache"
	"DistributedStorage/conf"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/pflag"
)

var cfg = pflag.StringP("conf", "c", "", "")

func main() {
	pflag.Parse()
	if err := conf.TestInit(*cfg); err != nil {
		panic(err)
	}
	cache.Init()

	err := cache.Set("SET", "name", "wrath", "EX", "1")
	m, err := redis.String(cache.Get("GET", "name"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(m)
}
