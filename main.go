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
	"tcpserver/configServer"
	"tcpserver/database/redisdatabase"
	"tcpserver/database/xormdatabase"
	"tcpserver/logServer"
	"tcpserver/server"
)

func main() {
	configServer.ParseConfig("./config/config.json")
	config := configServer.GetConfig()
	xormdatabase.InitXormEngine(config)
	redisdatabase.InitRedis(config)
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
