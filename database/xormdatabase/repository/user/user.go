//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: user.go
// description: 具体的数据库操作实现
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-05-15
//

package user

import (
	"github.com/InsideOfTheIndustry/TcpServe/database/xormdatabase"
	"github.com/InsideOfTheIndustry/TcpServe/logServer"
	"github.com/InsideOfTheIndustry/TcpServe/reposity"
)

// UserRepository 用户的dao操作
type UserRepository struct {
	*xormdatabase.XormEngine
}

// Create(user *entity.UserInfo) (int64,error) // 创建新用户 返回用户账号信息
func (ud UserRepository) SetFriend(launcher, accepter int64) (bool, error) {
	var friendinfo = UserFriend{
		Launcher: launcher,
		Accepter: accepter,
	}

	_, err := ud.Insert(&friendinfo)
	if err != nil {
		logServer.Error("添加好友关联失败:%s", err.Error())
		return false, err
	}

	return true, nil
}

// Query(useraccount int64) (*entity.UserInfo, error) // 查询用户信息
func (ud UserRepository) Query(useraccount int64) (*reposity.UserInfo, error) {
	var userinfo = UserInfo{}
	_, err := ud.Where("useraccount = ?", useraccount).Get(&userinfo)

	var userinfoentity = reposity.UserInfo{}

	if err != nil {
		logServer.Error("查询数据出错:(%s)", err.Error())
		return &userinfoentity, err
	}
	userinfoentity.UserAccount = userinfo.UserAccount
	userinfoentity.UserEmail = userinfo.UserEmail
	userinfoentity.UserName = userinfo.UserName
	userinfoentity.Signature = userinfo.Signature
	userinfoentity.Avatar = userinfo.Avatar
	userinfoentity.UserPassword = userinfo.UserPassword
	userinfoentity.UserAge = userinfo.UserAge
	userinfoentity.UserSex = userinfo.UserSex

	return &userinfoentity, nil
}

// Update(*entity.UserInfo) error // 更新用户信息 不包括头像信息
func (ud UserRepository) Update(user *reposity.UserInfo) error {
	var userindatabase = UserInfo{
		UserAccount:  user.UserAccount,
		UserName:     user.UserName,
		Signature:    user.Signature,
		UserPassword: user.UserPassword,
		UserAge:      user.UserAge,
		UserSex:      user.UserSex,
	}
	_, err := ud.Where("useraccount = ?", user.UserAccount).Update(userindatabase)
	if err != nil {
		logServer.Error("更新用户失败：（%s）", err.Error())
		return err
	}
	logServer.Error("更新用户成功。")
	return nil
}

// QueryFriends(useraccount int64)([]entity.FriendInfo, error) // 查询用户好友信息
func (ud UserRepository) QueryFriends(useraccount int64) (reposity.FriendInfo, error) {
	var friendlauchers = make([]UserFriend, 0)
	var friendaccepters = make([]UserFriend, 0)
	var friendsinfo = reposity.FriendInfo{
		UserAccount: useraccount,
	}

	// 查询朋友信息时 需要从发起者和接收者两处查询
	err := ud.Where("launcher = ?", useraccount).Find(&friendlauchers)
	if err != nil {
		logServer.Error("查询出现错误:(%s)", err.Error())
		return friendsinfo, err
	}
	err = ud.Where("accepter = ?", useraccount).Find(&friendaccepters)
	if err != nil {
		logServer.Error("查询出现错误:(%s)", err.Error())
		return friendsinfo, err
	}

	var friendsinfolist = make([]reposity.UserInfo, 0)
	// 查询朋友的具体信息
	for i := range friendlauchers {
		var friendinfo = UserInfo{}
		ud.Where("useraccount = ?", friendlauchers[i].Launcher).Get(&friendinfo)
		var entityuserinfo = reposity.UserInfo{
			UserAccount:  friendinfo.UserAccount,
			UserEmail:    friendinfo.UserEmail,
			UserName:     friendinfo.UserName,
			Signature:    friendinfo.Signature,
			Avatar:       friendinfo.Avatar,
			UserPassword: friendinfo.UserPassword,
			UserAge:      friendinfo.UserAge,
			UserSex:      friendinfo.UserSex,
		}
		friendsinfolist = append(friendsinfolist, entityuserinfo)
	}

	for i := range friendaccepters {
		var friendinfo = UserInfo{}
		ud.Where("useraccount = ?", friendaccepters[i].Accepter).Get(&friendinfo)
		var entityuserinfo = reposity.UserInfo{
			UserAccount:  friendinfo.UserAccount,
			UserEmail:    friendinfo.UserEmail,
			UserName:     friendinfo.UserName,
			Signature:    friendinfo.Signature,
			Avatar:       friendinfo.Avatar,
			UserPassword: friendinfo.UserPassword,
			UserAge:      friendinfo.UserAge,
			UserSex:      friendinfo.UserSex,
		}
		friendsinfolist = append(friendsinfolist, entityuserinfo)
	}
	friendsinfo.Friends = friendsinfolist
	return friendsinfo, nil

}

// QueryEmailIfAlreadyUse(email string) (bool, error)           // 查询邮箱是否已经注册
func (ud UserRepository) QueryEmailIfAlreadyUse(email string) (bool, error) {
	var userinfo = UserInfo{}
	count, err := ud.Where("useremail = ?", email).Count(userinfo)
	if err != nil {
		logServer.Error("查询邮箱是否被注册出现错误:(%s)", err.Error())
		return true, err
	}
	if count == 0 {
		return false, nil
	}
	return true, nil
}
