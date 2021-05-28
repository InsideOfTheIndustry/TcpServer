//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: config_struct.go
// description: 日数据结构定义
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-04-27
//

package configServer

// 配置文件结构
type ConfigStruct struct {
	Ip string `json:"ip"` // ip地址
	Port string `json:"port"` // 端口号
	ConfigPath string `json:"configpath"` // 日志路由
	Tcpserver  TcpServer `json:"tcpserver"` // tcp服务器配置
	Database Database `json:"database"` // 数据库配置文件
	Redis Redis `json:"redis"` // Redis缓存数据库
}

// tcp服务器ip和port
type TcpServer struct {
	Ip string `json:"ip"` // tcp服务器的ip地址
	Port string `json:"port"` // tcp服务器的端口
}

//Database 数据库配置文件
type Database struct {
	Type string `json:"type"` // 数据库类型 mysql..
	User string `json:"user"` // 用户名
	Password string `json:"password"` // 密码
	Host string `json:"host"` // IP
	Port string `json:"port"` // 端口
	DBName string `json:"dbName"` // 数据库名
	Charset string `json:"charset"` // 编码方式
	Showsql bool `json:"showsql"` // 是否显示数据库查询语句
}


// Redis缓存数据库
type Redis struct {
	Addr string `json:"addr"` // 数据库地址
	Password string `json:"password"` // 数据库连接密码
	Db int `json:"db"` // 使用的数据库
}

