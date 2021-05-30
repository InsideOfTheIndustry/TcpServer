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
	"log"

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
		log.Panic("读取配置文件失败")
		panic("读取配置文件失败")
	}

	// 服务参数
	cfgApplication = settingCfg.Sub("settings.application")
	if cfgApplication == nil {
		log.Panic("can not find application")
		panic("can not find application")
	}
	ApplicationConfig = InitApplication(cfgApplication)

	// 配置TCP服务器
	cfgTcpserver = settingCfg.Sub("settings.tcpserver")
	if cfgApplication == nil {
		log.Panic("can not find tcpserver")
		panic("can not find tcpserver")
	}
	TcpServerConfig = InitTcpServer(cfgTcpserver)

	// 数据库配置
	cfgDatabase = settingCfg.Sub("settings.database")
	if cfgDatabase == nil {
		log.Panic("can not find database")
		panic("can not find database")
	}
	DatabaseConfig = InitDatabase(cfgDatabase)
	// redis配置
	cfgRedis = settingCfg.Sub("settings.redis")
	if cfgDatabase == nil {
		log.Panic("can not find redis")
		panic("can not find redis")
	}
	RedisConfig = InitRedis(cfgRedis)
}
