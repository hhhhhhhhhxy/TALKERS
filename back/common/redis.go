package common

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var MyRedis *redis.Client

func RedisInit() *redis.Client {
	host := viper.GetString("redis.host")
	port := viper.GetString("redis.port")
	addr := fmt.Sprintf("%s:%s", host, port)
	rds := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})
	MyRedis = rds
	return rds
}
