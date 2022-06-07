package datasource

import (
	"go-shop/config"
	"go-shop/utils"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

// redis
func NewRedisConn() (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisConfig.Addr + ":" + config.RedisConfig.Port, // 指定
		Password: config.RedisConfig.Pwd,
		DB:       0, // redis一共16个库，指定其中一个库即可 使用default DB
	})

	//ping redis数据库
	_, err := rdb.Ping().Result()
	if err != nil {
		utils.Logger.Error("redis连接失败", zap.Any("err", err))
		return nil, err
	}

	utils.Logger.Info("redis连接成功")

	return rdb, nil
}
