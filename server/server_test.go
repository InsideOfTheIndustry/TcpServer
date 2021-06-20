package server

import (
	"sync"
	"testing"
)

// func TestNewtcp(t *testing.T) {

// 	logServer.Setup("info")               // 设置日志等级
// 	config.Setup("../config/config.yaml") // 读取配置文件
// 	// 初始化xorm引擎
// 	if err := xormdatabase.InitXormEngine(); err != nil {
// 		logServer.Error("初始化xorm引擎失败: %s", err.Error())
// 	}
// 	redisdatabase.InitRedis()                               // 初始化reds
// 	ctx, cancel := context.WithCancel(context.Background()) // 全局上下文控制
// 	// defer cancel()
// 	tcpserver, err := NewTcpServer(ctx) // 启动tcp服务
// 	if err != nil {
// 		t.Errorf("启动tcp失败：%s", err.Error())
// 	}
// 	var count = 0
// 	for count < 10 {
// 		time.Sleep(time.Millisecond * 1000)
// 		count += 1
// 	}

// 	cancel()
// 	<-ctx.Done()
// 	time.Sleep(time.Microsecond * 1000 * 50)
// 	for i := range tcpserver.conn {
// 		_, err = tcpserver.conn[i].Write([]byte("s"))
// 		if err == nil {
// 			t.Error("有连接未关闭")
// 		}
// 	}

// 	logServer.Info("Tcp服务停止")
// }

func TestTcpGroup(t *testing.T) {
	// logServer.Setup("info")               // 设置日志等级
	// config.Setup("../config/config.yaml") // 读取配置文件
	// // 初始化xorm引擎
	// if err := xormdatabase.InitXormEngine(); err != nil {
	// 	logServer.Error("初始化xorm引擎失败: %s", err.Error())
	// }
	// redis.InitRedis()                                       // 初始化reds
	// ctx, cancel := context.WithCancel(context.Background()) // 全局上下文控制
	// defer cancel()
	// tcpserver, err := NewTcpServer(ctx) // 启动tcp服务
	// if err != nil {
	// 	t.Errorf("启动tcp失败：%s", err.Error())
	// }
	// var count = 0
	// for count < 10 {
	// 	time.Sleep(time.Millisecond * 1000)
	// 	count += 1
	// }

	// // cancel()
	// <-ctx.Done()
	// time.Sleep(time.Microsecond * 1000 * 50)
	// for i := range tcpserver.conn {
	// 	_, err = tcpserver.conn[i].Write([]byte("s"))
	// 	if err == nil {
	// 		t.Error("有连接未关闭")
	// 	}
	// }

	// logServer.Info("Tcp服务停止")
	var maps = make(map[string]string)
	maps["1"] = "1"
	delete(maps, "2")
	var maps2 = sync.Map{}
	maps2.Store("111", 22)
	maps2.Delete("111")
	maps2.Delete("111")

}
