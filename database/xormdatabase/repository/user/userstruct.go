//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: userstruct.go
// description: 数据库内数据的结构
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-05-17
//
package user

import "time"

// UserInfo 用户实体
type UserInfo struct {
	UserAccount  int64  `xorm:"useraccount"`  // 用户账号
	UserEmail    string `xorm:"useremail"`    // 邮箱号
	UserName     string `xorm:"username"`     // 用户名
	Signature    string `xorm:"signature"`    // 用户个性签名
	Avatar       string `xorm:"avatar"`       // 用户头像
	UserPassword string `xorm:"userpassword"` // 用户密码
	UserAge      int64  `xorm:"userage"`      // 用户年龄
	UserSex      int64  `xorm:"usersex"`      // 用户性别
	Online       int8   `xorm:"online"`       // 用户在线状态
	Delete       int8   `xorm:"delete"`       // 用户是否被删除
}

// UserFriend 朋友间的相互联系
type UserFriend struct {
	Launcher int64 `xorm:"launcher"` // 好友发起者
	Accepter int64 `xorm:"accepter"` // 好友接受者
}

// GroupInfo 群聊信息
type GroupInfo struct {
	Groupid     int64     `xorm:"groupid"`     // 群聊号
	GroupName   string    `xorm:"groupname"`   // 群聊名称
	GroupIntro  string    `xorm:"groupintro"`  // 群聊简介
	GroupOwner  int64     `xorm:"groupowner"`  // 群主
	GroupAvatar string    `xorm:"groupavatar"` // 群头像
	Deleted     int8      `xorm:"delete"`      // 是否删除
	CreateAt    time.Time `xorm:"createat"`    // 创建时间
}

// UserGroup 用户所在的群
type UserGroup struct {
	Useraccount     int64  `xorm:"useraccount"`     // 用户账号
	Groupid         int64  `xorm:"groupid"`         // 群聊号
	UserNameInGroup string `xorm:"usernameingroup"` // 用户在群中的昵称
}

// GroupMemberInfo 群内用户的信息
type GroupMemberInfo struct {
	UserName    string `xorm:"username"`    // 用户名
	UserAccount int64  `xorm:"useraccount"` // 用户账号
}
