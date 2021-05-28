//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: initredis.go
// description: redis数据库实现
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-05-15
//

package redisdatabase

import (
	"tcpserver/configServer"
	"tcpserver/logServer"

	"github.com/go-redis/redis/v8"
)

// Redis引擎
type RedisEngine struct {
	*redis.Client
}

var RedisClient *RedisEngine

// InitRedis 初始化redis连接
func InitRedis(config *configServer.ConfigStruct) {
	redisconfig := config.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisconfig.Addr,
		Password: redisconfig.Password,
		DB:       redisconfig.Db,
	})

	var newredisclient = &RedisEngine{}
	newredisclient.Client = rdb

	RedisClient = newredisclient
	logServer.Info("redis数据库连接成功。")
}
