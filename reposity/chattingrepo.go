//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: chattingrepo.go
// description: 信息存储库
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-05-20
//

package reposity

// ChattingReposity 聊天存储库
type ChattingReposity interface {
	Query(useraccount int64) (*UserInfo, error)       // 用户是否存在
	SetFriend(launcher, accepter int64) (bool, error) // 建立朋友关系
}

// ChttingCacheReposity 聊天缓存库
type ChatingCacheReposity interface {
	GetToken(useraccount int64) (string, error) // token 是否存在是否与用户提供的相同
}
