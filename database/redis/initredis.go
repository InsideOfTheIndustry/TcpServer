//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: initredis.go
// description: redis数据库实现
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-05-15
//

package redis

import (
	"github.com/InsideOfTheIndustry/TcpServe/config"
	"github.com/InsideOfTheIndustry/TcpServe/logServer"

	"github.com/go-redis/redis/v8"
)

// Redis引擎
type RedisEngine struct {
	*redis.Client
}

var RedisClient *RedisEngine

// InitRedis 初始化redis连接
func InitRedis() {

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisConfig.Addr + ":" + config.RedisConfig.Port,
		Password: config.RedisConfig.Password,
		DB:       config.RedisConfig.Db,
	})

	var newredisclient = &RedisEngine{}
	newredisclient.Client = rdb

	RedisClient = newredisclient
	logServer.Info("redis数据库连接成功。")
}
