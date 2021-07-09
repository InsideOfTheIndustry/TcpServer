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
	"errors"
	"strconv"

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
	if err := ud.Where("launcher = ?", useraccount).Find(&friendlauchers); err != nil {
		logServer.Error("查询出现错误:(%s)", err.Error())
		return friendsinfo, err
	}

	if err := ud.Where("accepter = ?", useraccount).Find(&friendaccepters); err != nil {
		logServer.Error("查询出现错误:(%s)", err.Error())
		return friendsinfo, err
	}

	// 获取被定义为接受者的朋友
	var accepter = make([]int64, 0, len(friendlauchers))
	for i := range friendlauchers {
		accepter = append(accepter, friendlauchers[i].Accepter)
	}
	// 获取定义为发起者的朋友
	var launcher = make([]int64, 0, len(friendaccepters))
	for i := range friendaccepters {
		launcher = append(launcher, friendaccepters[i].Launcher)
	}

	var friendsinfolistaccepter = make([]reposity.UserInfo, 0, len(accepter))
	var friendsinfolistlauncher = make([]reposity.UserInfo, 0, len(launcher))

	if err := ud.In("useraccount", accepter).Find(&friendsinfolistaccepter); err != nil {
		logServer.Error("查询朋友信息出错:%s", err.Error())
		return friendsinfo, err
	}
	if err := ud.In("useraccount", launcher).Find(&friendsinfolistlauncher); err != nil {
		logServer.Error("查询朋友信息出错:%s", err.Error())
		return friendsinfo, err
	}
	friendsinfo.Friends = append(friendsinfo.Friends, friendsinfolistaccepter...)
	friendsinfo.Friends = append(friendsinfo.Friends, friendsinfolistlauncher...)
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

// QueryAllGroup 查询系统中的所有群聊
func (ud UserRepository) QueryAllGroup() ([]reposity.GroupInfo, error) {
	var grouplist = make([]GroupInfo, 0)

	if err := ud.Find(&grouplist); err != nil {
		logServer.Error("查询所有群聊信息失败:%s", err.Error())
		return make([]reposity.GroupInfo, 0), err
	}

	var groupinfolist = make([]reposity.GroupInfo, len(grouplist))

	for i := range grouplist {
		if grouplist[i].Deleted == 0 {
			groupinfolist[i].Groupid = grouplist[i].Groupid
			groupinfolist[i].GroupIntro = grouplist[i].GroupIntro
			groupinfolist[i].GroupName = grouplist[i].GroupName
		}

	}

	return groupinfolist, nil
}

// QueryGroupOfUser 查询用户所在的群
func (ud UserRepository) QueryGroupOfUser(useraccount int64) ([]string, error) {
	var usergroupinfo = make([]UserGroup, 0)
	if err := ud.Where("useraccount = ?", useraccount).Find(&usergroupinfo); err != nil {
		logServer.Error("查询用户所在的群失败:%s", err.Error())
		return make([]string, 0), err
	}

	var groupidlist = make([]string, len(usergroupinfo))
	for i := range usergroupinfo {
		groupid := strconv.FormatInt(usergroupinfo[i].Groupid, 10)
		groupidlist = append(groupidlist, groupid)
	}
	return groupidlist, nil

}

// UpdateUserOnlineStatue 更新用户在线状态
func (ud UserRepository) UpdateUserOnlineStatue(useraccount int64, status bool) error {
	var online = 0
	if status {
		online = 1
	}
	var us = UserInfo{
		UserAccount: useraccount,
		Online:      int8(online),
	}
	if _, err := ud.Where("useraccount = ?", useraccount).Cols("online").Update(us); err != nil {
		logServer.Error("更新用户状态失败:%s", err.Error())
		return err
	}

	return nil

}

// QueryGroupMembers 查询群聊用户
func (ud UserRepository) QueryGroupMembers(groupid int64) ([]reposity.GroupMemberInfo, error) {
	var groupmembers = make([]GroupMemberInfo, 0)

	if err := ud.Sql("SELECT username,UserGroup.useraccount FROM UserInfo , UserGroup   WHERE  UserGroup.groupid = ? AND UserInfo.useraccount = UserGroup.useraccount", groupid).Find(&groupmembers); err != nil {
		return []reposity.GroupMemberInfo{}, err
	}

	gm := make([]reposity.GroupMemberInfo, len(groupmembers))

	for i := range groupmembers {
		gm[i] = reposity.GroupMemberInfo(groupmembers[i])
	}

	return gm, nil
}

// QueryGroupMembersCount 查询群聊用户数量
func (ud UserRepository) QueryGroupMembersCount(groupid int64) (int64, error) {

	numbers, err := ud.Table("UserGroup").Where("groupid = ?", groupid).Count()
	if err != nil {
		return 0, err
	}

	return numbers, nil
}

// QueryGroupInfo 查询群聊信息
func (ud UserRepository) QueryGroupInfo(groupid int64) (reposity.GroupInfo, error) {
	var groupinfo = GroupInfo{Groupid: groupid}
	ok, err := ud.Get(&groupinfo)
	if err != nil {
		logServer.Error("查询群聊信息失败:%s", err.Error())
		return reposity.GroupInfo{}, err
	}
	logServer.Info("ok:%v,value:%v", ok, groupinfo)
	if groupinfo.Deleted == 1 {
		return reposity.GroupInfo{}, nil
	}

	var groupinforeturn = reposity.GroupInfo{
		Groupid:     groupinfo.Groupid,
		GroupName:   groupinfo.GroupName,
		GroupIntro:  groupinfo.GroupIntro,
		GroupAvatar: groupinfo.GroupAvatar,
		GroupOwner:  groupinfo.GroupOwner,
		CreateAt:    groupinfo.CreateAt,
		Deleted:     groupinfo.Deleted,
	}

	return groupinforeturn, nil
}

// QueryIfUserInGroup 查询用户是否在群内
func (ud UserRepository) QueryIfUserInGroup(useraccount int64, groupid int64) (bool, error) {
	var usergroupinfo = UserGroup{}
	ifexist, err := ud.Where("useraccount = ?", useraccount).And("groupid = ?", groupid).Get(&usergroupinfo)
	if err != nil {
		logServer.Error("查询用户所在的群失败:%s", err.Error())
		return false, err
	}

	return ifexist, nil
}

// AddUserToGroup 将用户加入群聊
func (ud UserRepository) AddUserToGroup(useraccount, groupid int64) error {

	var userinfo = UserInfo{}
	_, err := ud.Where("useraccount = ?", useraccount).Get(&userinfo)
	if err != nil {
		logServer.Error("查询用户失败:%s", err.Error())
		return err
	}

	if userinfo.UserAccount == 0 {

		return errors.New("用户不存在!")
	}

	// 插入用户 群聊关系 准备加一个时间字段
	var usergroup = UserGroup{
		Useraccount:     useraccount,
		Groupid:         groupid,
		UserNameInGroup: userinfo.UserName,
	}
	if _, err := ud.InsertOne(usergroup); err != nil {
		logServer.Error("将用户加入默认群聊失败:%s", err.Error())
		return err
	}

	return nil
}
