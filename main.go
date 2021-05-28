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
	config.Setup("config/settings.yaml")

	xormdatabase.InitXormEngine()
	redisdatabase.InitRedis()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	logServer.SetFileLevel("info")
	go server.NewTcpServer(ctx)
	for {
		select {
		case <-ctx.Done():
			goto stopposition
		}
	}
stopposition:
	logServer.Info("Tcp服务停止")

}
