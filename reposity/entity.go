//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: entity.go
// description: 实体库
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-05-20
//

package reposity

// 用户实体
type UserInfo struct {
	UserAccount  int64  // 用户账号
	UserEmail    string // 邮箱号
	UserName     string // 用户名
	Signature    string // 用户个性签名
	Avatar       string // 用户头像
	UserPassword string // 用户密码
	UserAge      int64  // 用户年龄
	UserSex      int64  // 用户性别
}

// 朋友间的相互联系
type UserFriend struct {
	Launcher int64 // 好友发起者
	Accepter int64 // 好友接受者
}

// 用户朋友信息
type FriendInfo struct {
	UserAccount int64      // 用户名
	Friends     []UserInfo // 好友信息
}
