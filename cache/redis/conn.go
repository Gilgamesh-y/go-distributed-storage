package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
	"time"
)

var pool *redis.Pool

func Init() {
	pool = &redis.Pool{
		MaxIdle: 50,
		MaxActive: 30,
		IdleTimeout: 300*time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", viper.GetString("cache_host") + ":" + viper.GetString("cache_port"))
			if err != nil {
				fmt.Println(err)
				return nil, err
			}
			if viper.GetString("cache_pass") != "" {
				if _, err = c.Do("AUTH", viper.GetString("cache_pass")); err != nil {
					c.Close()
					return nil, err
				}
			}

			return c, nil
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

func GetPool() *redis.Pool {
	return pool
}

func Set(action string, key string, args ...interface{}) error {
	conn := GetPool().Get()
	defer conn.Close()
	_, err := conn.Do(action, key, args)
	return err
}

func Get(action string, key string) (interface{}, error) {
	conn := GetPool().Get()
	defer conn.Close()
	return conn.Do(action, key)
}