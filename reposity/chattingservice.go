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
	if userinfo.UserEmail == "" || userinfo.Delete == 1 {
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

// UpdateUserOnlineStatus 更新用户在线状态
func (us UserService) UpdateUserOnlineStatus(useraccount int64, status bool) error {
	return us.ChattingReposity.UpdateUserOnlineStatue(useraccount, status)
}

// QueryGroupMembers 查询群聊成员信息
func (us UserService) QueryGroupMembers(groupid int64) ([]GroupMemberInfo, error) {
	return us.ChattingReposity.QueryGroupMembers(groupid)
}

// QueryGroupInfo 查询群聊信息
func (us UserService) QueryGroupInfo(groupid int64) (GroupInfo, error) {
	return us.ChattingReposity.QueryGroupInfo(groupid)
}

// QueryGroupOfUser 查询用户拥有的群聊
func (us UserService) QueryGroupOfUser(groupid int64) ([]string, error) {
	return us.ChattingReposity.QueryGroupOfUser(groupid)
}

//QueryIfUserInGroup 查询用户是否再群内
func (us UserService) QueryIfUserInGroup(useraccount int64, groupid int64) (bool, error) {
	return us.ChattingReposity.QueryIfUserInGroup(useraccount, groupid)
}

//  AddUserToGroup 将用户加入群聊
func (us UserService) AddUserToGroup(useraccount, groupid int64) error {
	return us.ChattingReposity.AddUserToGroup(useraccount, groupid)
}

// QueryGroupMembersCount 查询群内用户数
func (us UserService) QueryGroupMembersCount(groupid int64) (int64, error) {
	return us.ChattingReposity.QueryGroupMembersCount(groupid)
}
