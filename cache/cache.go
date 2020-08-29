package cache

import (
	credis "DistributedStorage/cache/redis"
	"github.com/spf13/viper"
)

var cacheDriver string

func Init() {
	credis.Init()
	cacheDriver = viper.GetString("cache_driver")
}

func Set(key string, args ...interface{}) error {
	switch cacheDriver {
	case "redis":
		return credis.Set(key, args)
	default:
		panic("未设置缓存驱动")
	}
}

func GetString(key string) (string, error) {
	switch cacheDriver {
		case "redis":
			return credis.Get(key)
		default:
			panic("未设置缓存驱动")
	}
}