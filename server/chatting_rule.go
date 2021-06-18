//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: chattingRule.go
// description: 通信规则
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-04-28
//

package server

import "time"

type MessageTypes int

const (
	HeartBeat                 MessageTypes = 0   // 心跳信息
	FirstConnect              MessageTypes = 1   // 初次连接
	SendMessage               MessageTypes = 2   // 发送信息
	SendGroupMessage          MessageTypes = 22  // 发送群聊信息
	CloseConnect              MessageTypes = 4   // 断开连接
	UserNotOnline             MessageTypes = 404 // 用户不在线
	SendFriendRequest         MessageTypes = 3   // 发送交友请求
	AcceptFriendRequest       MessageTypes = 33  // 接收好友请求
	RejectFriendRequest       MessageTypes = 333 // 拒绝好友请求
	SendInfoSuccess           MessageTypes = 200 // 发出聊天信息成功
	SendInfoFaild             MessageTypes = 400 // 发出聊天信息失败
	SendGroupInfoSuccess      MessageTypes = 220 // 发出群聊信息成功
	SendGroupInfoFaild        MessageTypes = 440 // 发出群聊信息失败
	FriendMakeInfoSendSuccess MessageTypes = 201 // 发出好友相关的请求的信成功
	FriendMakeInfoSendFail    MessageTypes = 402 // 发出好友相关的请求的信失败
	AuthorizationFail         MessageTypes = 500 // 验证token失败
	OnlineStatus              MessageTypes = 222 // 上线了
	NotOlineStatus            MessageTypes = 444 // 离线了
	OtherPlaceLogin           MessageTypes = 88  // 在其他地方登录
)

// Message 信息传递结构
type Message struct {
	MessageType MessageTypes `json:"messagetype"` // 消息的类型 比如心跳、发送信息等
	Token       string       `json:"token"`       // token用于验证用户是否登录
	Message     string       `json:"message"`     // 发送的信息
	Sender      string       `json:"sender"`      // 发送者账号
	Receiver    string       `json:"receiver"`    // 接收者账号
	Groupid     string       `json:"groupid"`     // 群聊id
	SendTime    time.Time    `json:"time"`        // 时间信息
}
