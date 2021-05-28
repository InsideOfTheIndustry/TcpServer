//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: user_test.go
// description: 用户表单元测试
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-05-15
//

package user

import (
	"testing"

	"github.com/InsideOfTheIndustry/tcpserver/reposity"

	"github.com/InsideOfTheIndustry/TcpServe/database/xormdatabase"
	"github.com/InsideOfTheIndustry/TcpServe/logServer"
)

func TestCreate(t *testing.T) {
	// 读取配置文件
	configServer.ParseConfig("../../../../config/config.json")
	config := configServer.GetConfig()
	// 测试数据库
	var userdao = UserRepository{}
	xormdatabase.InitXormEngine(config)
	userdao.XormEngine = xormdatabase.DBEngine
	var userinfo = reposity.UserInfo{
		UserEmail:    "xxx2",
		UserName:     "小猪猪",
		UserPassword: "123456",
		Signature:    "xxxxx",
		UserAge:      19,
		UserSex:      1,
		UserAccount:  100000,
	}
	// account, _ := userdao.Create(&userinfo)
	// users, err := userdao.Query(12313131)
	ifre, err := userdao.QueryEmailIfAlreadyUse("12138")
	friends, _ := userdao.QueryFriends(12138)
	userdao.Update(&userinfo)

	//logServer.Info("用户账号为:%v", account)
	//logServer.Info("用户是否存在:%v,%v", *users, err)
	logServer.Info("邮箱是否已被注册:%v", ifre)
	logServer.Info("错误信息:%v", err)
	logServer.Info("好友信息:%v", friends)
}
