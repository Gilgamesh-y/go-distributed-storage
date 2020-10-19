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

func Set(action string, args ...interface{}) error {
	switch cacheDriver {
	case "redis":
		return credis.Set(action, args...)
	default:
		panic("未设置缓存驱动")
	}
}

func Get(action string, key string) (interface{}, error) {
	switch cacheDriver {
		case "redis":
			return credis.Get(action, key)
		default:
			panic("未设置缓存驱动")
	}
}