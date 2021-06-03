//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: chattingservice.go
// description: 用户服务
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-05-20
//
package reposity

import (
	"time"

	"github.com/InsideOfTheIndustry/TcpServe/logServer"
	"github.com/InsideOfTheIndustry/TcpServe/utils/jwt"
)

// UserService 用户领域服务
type UserService struct {
	ChattingReposity     // 聊天基本存储库
	ChatingCacheReposity // 缓存库
}

// IfExistUser 判断用户是否存在
func (us UserService) IfExistUser(useraccount int64) (bool, error) {
	userinfo, err := us.ChattingReposity.Query(useraccount)
	if err != nil {
		logServer.Error("查询用户出现错误:%s", err.Error())
		return false, err
	}
	if userinfo.UserEmail == "" {
		return false, nil
	}
	return true, nil
}

// IfTokenSame判断token是否存在且相同
func (us UserService) IfTokenSameAndNotExpired(useraccount int64, token string) (bool, error) {
	tokencache, err := us.ChatingCacheReposity.GetToken(useraccount)
	if err != nil {
		logServer.Error("读取缓存失败:%s", err.Error())
		return false, err
	}
	if tokencache != token {
		return false, nil
	}
	claim, err := jwt.ParseToken(token)
	if err != nil {
		logServer.Error("token解析失败:%s", err.Error())
		return false, err
	}
	if claim.ExpiresAt < time.Now().Unix() {
		return false, nil
	}
	return true, nil
}

// BuildFriend 建立好友关系
func (us UserService) BuildFriend(launcher, accepter int64) (bool, error) {
	ifsuccess, err := us.ChattingReposity.SetFriend(launcher, accepter)
	return ifsuccess, err
}
