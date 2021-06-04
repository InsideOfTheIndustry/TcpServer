//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: main.go
// description: 具体的main函数
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-04-27
//

package main

import (
	"context"

	"github.com/InsideOfTheIndustry/TcpServe/config"
	redisdatabase "github.com/InsideOfTheIndustry/TcpServe/database/redis"
	"github.com/InsideOfTheIndustry/TcpServe/database/xormdatabase"
	"github.com/InsideOfTheIndustry/TcpServe/logServer"
	"github.com/InsideOfTheIndustry/TcpServe/server"
)

func main() {
	logServer.Setup("info")              // 设置日志等级
	config.Setup("./config/config.yaml") // 读取配置文件

	// 初始化xorm引擎
	if err := xormdatabase.InitXormEngine(); err != nil {
		logServer.Error("初始化xorm引擎失败: %s", err.Error())
	}
	redisdatabase.InitRedis()                               // 初始化reds
	ctx, cancel := context.WithCancel(context.Background()) // 全局上下文控制
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	go server.NewTcpServer(ctx) // 启动tcp服务

	<-ctx.Done()

	logServer.Info("Tcp服务停止")

}
