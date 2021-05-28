//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: userstruct.go
// description: 数据库内数据的结构
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-05-17
//
package user

// 用户实体
type UserInfo struct {
	UserAccount  int64  `xorm:"useraccount"`  // 用户账号
	UserEmail    string `xorm:"useremail"`    // 邮箱号
	UserName     string `xorm:"username"`     // 用户名
	Signature    string `xorm:"signature"`    // 用户个性签名
	Avatar       string `xorm:"avatar"`       // 用户头像
	UserPassword string `xorm:"userpassword"` // 用户密码
	UserAge      int64  `xorm:"userage"`      // 用户年龄
	UserSex      int64  `xorm:"usersex"`      // 用户性别
}

// 朋友间的相互联系
type UserFriend struct {
	Launcher int64 `xorm:"launcher"` // 好友发起者
	Accepter int64 `xorm:"accepter"` // 好友接受者
}
