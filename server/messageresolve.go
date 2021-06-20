//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: messageresolve.go
// description: 信息类型对应的处理方式
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-05-20
//

package server

import (
	"encoding/json"
	"net"
	"strconv"
	"time"

	"github.com/InsideOfTheIndustry/TcpServe/logServer"
	"github.com/InsideOfTheIndustry/TcpServe/reposity"
)

// NewUserLoginIn 新连接加入
func (tcpserver *TcpServer) NewUserLoginIn(service reposity.UserService, useraccount int64, receiveMessage Message, connectidentify *ConnectIdentify) {

	// 检查是否存在用户
	exist, err := service.IfExistUser(useraccount)
	if err != nil || !exist {
		connectidentify.connect.Close()
		return
	}
	// 检查当前是否在线
	receiverConnOther, ok := tcpserver.connectionpool.Load(receiveMessage.Sender)

	if ok {
		logServer.Info("重复登录了")
		var receiverConn = receiverConnOther.(*net.TCPConn)
		SendCommonMessage(receiverConn, "tcpserver provider", receiveMessage.Sender, "您的账号被人登录了，您已下线。", "", OtherPlaceLogin)
		err := receiverConn.Close()
		if err != nil {
			logServer.Error("关闭连接失败:(%s)")
		}
	}

	// 如果重复登录，需要对其进行切换
	tcpserver.connectionpool.Store(receiveMessage.Sender, connectidentify.connect) // 将连接加入连接池
	connectidentify.ifinconnectpool = true
	connectidentify.expireat.Reset(30 * time.Second)
	connectidentify.useraccount = receiveMessage.Sender
	logServer.Info("用户：(%s)加入了聊天。", receiveMessage.Sender)

	// 需要将其加入对应群聊
	sender, _ := strconv.ParseInt(receiveMessage.Sender, 10, 64)
	chattinglist, err := tcpserver.service.QueryGroupOfUser(sender)
	if err != nil {
		return
	}

	for i := range chattinglist {
		groupi, ok := tcpserver.groupchatting.Load(chattinglist[i])
		if !ok {
			continue
		}
		group := groupi.(Group)
		group.lock.Lock()
		group.groupmember[receiveMessage.Sender] = struct{}{}
		group.lock.Unlock()
		tcpserver.groupchatting.Store(chattinglist[i], group)
	}

	// 更新在线信息
	if err := tcpserver.service.UpdateUserOnlineStatus(useraccount, true); err != nil {
		logServer.Error("更新状态失败:%s", err.Error())
	}

	// 向好友广播在线信息
	userfriend, _ := tcpserver.service.QueryFriends(useraccount)
	useraccounts := strconv.FormatInt(useraccount, 10)
	for i := range userfriend.Friends {
		friendaccount := strconv.FormatInt(userfriend.Friends[i].UserAccount, 10)
		conni, ok := tcpserver.connectionpool.Load(friendaccount)
		if ok {
			conn := conni.(*net.TCPConn)
			SendCommonMessage(conn, useraccounts, friendaccount, "i am online", "", OnlineStatus)
		}
	}

	//TODO: 若添加离线信息的话 需要先从数据库读取数据
}

// SendMessageToReceiver 发送信息给接收者 case:SendMessage
func (tcpserver *TcpServer) SendMessageToReceiver(receiveMessage Message, conn *net.TCPConn, successStatus, FailStatus MessageTypes) bool {

	// 从连接池中查找连接
	connectInterface, ok := tcpserver.connectionpool.Load(receiveMessage.Receiver)

	if ok {
		receiverConn := connectInterface.(*net.TCPConn)
		logServer.Info("用户：（%s）发送信息给用户：(%s)", receiveMessage.Sender, receiveMessage.Receiver)
		sendbyte, _ := json.Marshal(receiveMessage)
		_, err := receiverConn.Write(sendbyte)

		// 上线过，但是后续掉线了
		if err != nil {
			tcpserver.connectionpool.Delete(receiveMessage.Receiver)
			// TODO: 可以先存储到数据库内
			logServer.Error("信息发送失败:(%s)", err.Error())
			logServer.Info("用户:(%s)不在线", receiveMessage.Receiver)
			// 回复一个发送失败的信息
			SendCommonMessage(conn, receiveMessage.Receiver, receiveMessage.Sender, receiveMessage.Message, "", FailStatus)
			return false
		} else {
			SendCommonMessage(conn, receiveMessage.Receiver, receiveMessage.Sender, receiveMessage.Message, "", successStatus)
			logServer.Info("信息成功发送.")
			return true
		}
	} else {
		logServer.Info("用户:(%s)不在线", receiveMessage.Receiver)
		// 先回复当前用户
		// 需要判断是否存在此用户
		// TODO: 可以先存储到数据库内
		SendCommonMessage(conn, receiveMessage.Receiver, receiveMessage.Sender, receiveMessage.Message, "", FailStatus)
		return false
	}
}

// HeartBeatMessage 心跳信息
func (tcpserver *TcpServer) HeartBeatMessage(receiveMessage Message, connectidentify *ConnectIdentify) {
	logServer.Info("接收到心跳信息...")
	if receiveMessage.Sender != connectidentify.useraccount {
		logServer.Info("心跳发送者account不匹配，期望:%s, 实际发送者为:%s", connectidentify.useraccount, receiveMessage.Sender)
		if err := Gracefulclose(connectidentify.connect, "您的心跳异常"); err != nil {
			logServer.Error("心跳连接异常,连接关闭失败:%s", err.Error())
			return
		}

		tcpserver.connectionpool.Delete(receiveMessage.Sender)
		logServer.Info("用户：（%s）断开连接", receiveMessage.Sender)
	}

	connectidentify.expireat.Reset(30 * time.Second)

}

// CloseConnect 关闭连接做的事
func (tcpserver *TcpServer) CloseConnect(receiveMessage Message, conn *net.TCPConn, message string) {
	senderconninterface, ok := tcpserver.connectionpool.Load(receiveMessage.Sender)
	if ok {
		senderconn := senderconninterface.(*net.TCPConn)

		// gracefule close 断开前向对方发送通知
		if err := Gracefulclose(senderconn, message); err != nil {
			logServer.Error("关闭连接失败:%s", err.Error())
			return
		}

		tcpserver.connectionpool.Delete(receiveMessage.Sender)
		logServer.Info("用户：（%s）断开连接", receiveMessage.Sender)
	} else {
		conn.Close()
	}
}

// LaunchFrienRequest 发起好友请求
func (tcpserver *TcpServer) LaunchFriendRequest(receiveMessage Message, conn *net.TCPConn) {
	launcherint, _ := strconv.ParseInt(receiveMessage.Sender, 10, 64)
	accepterint, _ := strconv.ParseInt(receiveMessage.Receiver, 10, 64)
	friends, err := tcpserver.service.QueryFriends(launcherint)

	if err != nil {
		logServer.Error("查询用户好友信息出现错误:%s", err.Error())
		SendCommonMessage(conn, "tcpserver provider", receiveMessage.Sender, "查询好友失败", "", FriendMakeInfoSendFail)
		return
	}

	for i := range friends.Friends {
		if friends.Friends[i].UserAccount == accepterint {
			SendCommonMessage(conn, "tcpserver provider", receiveMessage.Sender, "找茬？", "", FriendMakeInfoSendFail)
			return
		}
	}

	// 目前只支持在线添加好友 首先将一个请求添加进好友交友队列
	var friendmake = friendMakeInfo{
		launcher:   launcherint,
		accepter:   accepterint,
		randomcode: receiveMessage.Sender + receiveMessage.Receiver,
	}

	// 发送信息给好友请求接收者
	var message = Message{
		MessageType: SendFriendRequest,
		Token:       "",
		Message:     friendmake.randomcode,
		Sender:      receiveMessage.Sender,
		Receiver:    receiveMessage.Receiver,
		SendTime:    receiveMessage.SendTime,
	}
	// 返回发送情况
	success := tcpserver.SendMessageToReceiver(message, conn, FriendMakeInfoSendSuccess, FriendMakeInfoSendFail)
	if success {
		tcpserver.friendMakeList.Store(friendmake.randomcode, friendmake)
	}
}

// AcceptFrienRequest 接受好友请求
func (tcpserver *TcpServer) AcceptFriendRequest(service reposity.UserService, receiveMessage Message, conn *net.TCPConn) {
	accepterint, _ := strconv.ParseInt(receiveMessage.Sender, 10, 64)
	launcherint, _ := strconv.ParseInt(receiveMessage.Receiver, 10, 64)

	var friendmakeinfopo = -1
	logServer.Info("收到接受好友请求信息。")

	value, ok := tcpserver.friendMakeList.Load(receiveMessage.Message)
	if ok {
		friendmakeinfo := value.(friendMakeInfo)
		if friendmakeinfo.randomcode == receiveMessage.Message && friendmakeinfo.launcher == launcherint && friendmakeinfo.accepter == accepterint {
			success, _ := service.ChattingReposity.SetFriend(launcherint, accepterint)
			if success {
				var message = Message{
					MessageType: AcceptFriendRequest,
					Token:       "",
					Message:     "",
					Sender:      receiveMessage.Sender,
					Receiver:    receiveMessage.Receiver,
					SendTime:    receiveMessage.SendTime,
				}
				// 返回发送情况
				tcpserver.SendMessageToReceiver(message, conn, FriendMakeInfoSendSuccess, FriendMakeInfoSendFail)
				friendmakeinfopo = 1
			}
		}
	}

	// 对好友添加列表进行删除
	if friendmakeinfopo != -1 {
		tcpserver.friendMakeList.Delete(receiveMessage.Message)
	}
}

// RejectFrienRequest 拒绝好友请求
func (tcpserver *TcpServer) RejectFrienRequest(receiveMessage Message, conn *net.TCPConn) {
	accepterint, _ := strconv.ParseInt(receiveMessage.Sender, 10, 64)
	launcherint, _ := strconv.ParseInt(receiveMessage.Receiver, 10, 64)

	var friendmakeinfopo = -1
	value, ok := tcpserver.friendMakeList.Load(receiveMessage.Message)
	if ok {
		friendmakeinfo := value.(friendMakeInfo)
		if friendmakeinfo.randomcode == receiveMessage.Message && friendmakeinfo.launcher == launcherint && friendmakeinfo.accepter == accepterint {

			var message = Message{
				MessageType: RejectFriendRequest,
				Token:       "",
				Message:     "",
				Sender:      receiveMessage.Sender,
				Receiver:    receiveMessage.Receiver,
				SendTime:    receiveMessage.SendTime,
			}
			// 返回发送情况
			tcpserver.SendMessageToReceiver(message, conn, FriendMakeInfoSendSuccess, FriendMakeInfoSendFail)
			friendmakeinfopo = 1

		}
	}

	// 对好友添加列表进行删除
	if friendmakeinfopo != -1 {
		tcpserver.friendMakeList.Delete(receiveMessage.Message)
	}
}

// Gracefulclose 优雅断开连接 断开前发送信息
func Gracefulclose(conn *net.TCPConn, message string) error {
	sendclosemessage := Message{
		MessageType: CloseConnect,
		Message:     message,
	}
	sendmessagebyte, err := json.Marshal(sendclosemessage)
	if err != nil {
		logServer.Error("反序列化json数据失败:%s", err.Error())
		if err := conn.Close(); err != nil {
			logServer.Error("关闭连接出现错误(%s)", err.Error())
			return err
		}
		return err
	}

	if _, err := conn.Write(sendmessagebyte); err != nil {
		logServer.Error("发送信息失败:%s", err.Error())
	}
	if err := conn.Close(); err != nil {
		logServer.Error("关闭连接出现错误(%s)", err.Error())
		return err
	}

	return nil
}
