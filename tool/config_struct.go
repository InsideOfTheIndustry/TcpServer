//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: config_struct.go
// description: 日数据结构定义
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-04-27
//

package tool

// 配置文件结构
type ConfigStruct struct {
	Ip         string `json:"Ip"`         // ip地址
	Port       string `json:"Port"`       // 端口号
	ConfigPath string `json:"ConfigPath"` // 日志路由
}

var _config *ConfigStruct
