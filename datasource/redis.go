package datasource

import (
	"go-shop/config"

	"github.com/kataras/iris/v12/sessions/sessiondb/redis"
)

/**
 * 返回Redis实例
 */
func NewRedis() *redis.Database {
	var database *redis.Database
	database = redis.New(redis.Config{
		Network:   config.RedisConfig.NetWork,
		Addr:      config.RedisConfig.Addr + ":" + config.RedisConfig.Port,
		Password:  config.RedisConfig.Pwd,
		Database:  "",
		MaxActive: 10,
		Timeout:   redis.DefaultRedisTimeout,
		Prefix:    config.RedisConfig.Prefix,
	})
	return database
}
