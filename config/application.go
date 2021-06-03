//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: config_struct.go
// description: 日数据结构定义
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-04-27
//

package config

import "github.com/spf13/viper"

// 配置文件结构
type Application struct {
	Ip         string
	Port       string
	ConfigPath string
	JwtKey     string
}

func InitApplication(cfg *viper.Viper) *Application {
	return &Application{
		Ip:         cfg.GetString("ip"),
		Port:       cfg.GetString("port"),
		ConfigPath: cfg.GetString("configpath"),
		JwtKey:     cfg.GetString("jwtkey"),
	}
}

var ApplicationConfig = new(Application)

// tcp服务器ip和port
type TcpServer struct {
	Ip   string
	Port string
}

func InitTcpServer(cfg *viper.Viper) *TcpServer {
	return &TcpServer{
		Ip:   cfg.GetString("ip"),
		Port: cfg.GetString("port"),
	}
}

var TcpServerConfig = new(TcpServer)
