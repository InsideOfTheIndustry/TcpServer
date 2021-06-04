package server

import (
	"context"
	"testing"
	"time"

	"github.com/InsideOfTheIndustry/TcpServe/config"
	redisdatabase "github.com/InsideOfTheIndustry/TcpServe/database/redis"
	"github.com/InsideOfTheIndustry/TcpServe/database/xormdatabase"
	"github.com/InsideOfTheIndustry/TcpServe/logServer"
)

func TestNewtcp(t *testing.T) {

	logServer.Setup("info")               // 设置日志等级
	config.Setup("../config/config.yaml") // 读取配置文件
	// 初始化xorm引擎
	if err := xormdatabase.InitXormEngine(); err != nil {
		logServer.Error("初始化xorm引擎失败: %s", err.Error())
	}
	redisdatabase.InitRedis()                               // 初始化reds
	ctx, cancel := context.WithCancel(context.Background()) // 全局上下文控制
	// defer cancel()
	go NewTcpServer(ctx) // 启动tcp服务

	var count = 0
	for count < 10 {
		time.Sleep(time.Millisecond * 1000)
		count += 1
	}

	cancel()
	<-ctx.Done()

	logServer.Info("Tcp服务停止")
	for {
	}
}
