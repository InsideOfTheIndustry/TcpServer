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
	Query(useraccount int64) (*UserInfo, error)                        // 用户是否存在
	SetFriend(launcher, accepter int64) (bool, error)                  // 建立朋友关系
	QueryAllGroup() ([]GroupInfo, error)                               // 查询所有的群聊信息
	QueryGroupOfUser(useraccount int64) ([]string, error)              // 查询用户所在的群
	UpdateUserOnlineStatue(useraccount int64, status bool) error       // 更新用户在线状态
	QueryFriends(useraccount int64) (FriendInfo, error)                // 查询好友信息
	QueryGroupMembers(groupid int64) ([]GroupMemberInfo, error)        // 查询群内用户信息
	QueryGroupInfo(groupid int64) (GroupInfo, error)                   // 查询群信息
	QueryIfUserInGroup(useraccount int64, groupid int64) (bool, error) // 查询群内是否存在该用户
	AddUserToGroup(useraccount, groupid int64) error                   // 将用户加入群聊
	QueryGroupMembersCount(groupid int64) (int64, error)               // 查询群内用户数
}

// ChttingCacheReposity 聊天缓存库
type ChatingCacheReposity interface {
	GetToken(useraccount int64) (string, error) // token 是否存在是否与用户提供的相同
}
