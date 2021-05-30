//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: config.go
// description: 配置文件读取
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-04-27
//

package config

import (
	"github.com/InsideOfTheIndustry/TcpServe/logServer"
	"github.com/spf13/viper"
)

var cfgApplication *viper.Viper
var cfgTcpserver *viper.Viper
var cfgDatabase *viper.Viper
var cfgRedis *viper.Viper

func Setup(path string) {
	settingCfg := viper.New()
	settingCfg.SetConfigFile(path)
	if err := settingCfg.ReadInConfig(); err != nil {
		logServer.Error("配置文件读取失败:%s", err.Error())
		return
	}

	// 服务参数
	cfgApplication = settingCfg.Sub("settings.application")
	if cfgApplication == nil {
		logServer.Error("找不到application设置")
		return
	}
	ApplicationConfig = InitApplication(cfgApplication)

	// 配置TCP服务器
	cfgTcpserver = settingCfg.Sub("settings.tcpserver")
	if cfgTcpserver == nil {
		logServer.Error("找不到tcpserver设置")
		return
	}
	TcpServerConfig = InitTcpServer(cfgTcpserver)

	// 数据库配置
	cfgDatabase = settingCfg.Sub("settings.database")
	if cfgDatabase == nil {
		logServer.Error("找不到database设置")
		return
	}
	DatabaseConfig = InitDatabase(cfgDatabase)
	// redis配置
	cfgRedis = settingCfg.Sub("settings.redis")
	if cfgRedis == nil {
		logServer.Error("找不到redis设置")
		return
	}
	RedisConfig = InitRedis(cfgRedis)
}
