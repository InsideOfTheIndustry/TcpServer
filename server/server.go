//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: server.go
// description: tcp服务端
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-04-27
//

package server

import (
	"context"
	"encoding/json"
	"net"
	"strconv"
	"sync"
	"time"

	redisdatabase "github.com/InsideOfTheIndustry/TcpServe/database/redis"
	redisuser "github.com/InsideOfTheIndustry/TcpServe/database/redis/user"
	"github.com/InsideOfTheIndustry/TcpServe/database/xormdatabase"
	xormuser "github.com/InsideOfTheIndustry/TcpServe/database/xormdatabase/repository/user"
	"github.com/InsideOfTheIndustry/TcpServe/logServer"
	"github.com/InsideOfTheIndustry/TcpServe/reposity"
)

// Tcp服务器
type TcpServer struct {
	addr             net.TCPAddr          // 服务器地址
	connectionpool   sync.Map             // 用户连接池
	listener         net.TCPListener      // tcp监听器
	ctx              context.Context      // 上下文
	cancel           context.CancelFunc   // 退出回调
	friendMakeList   sync.Map             // 交友信息用于存储特定的聊天码 后续将采用键值对 来进行优化
	groupchatting    sync.Map             // 群聊
	groupmessagechan chan Message         // 群聊消息管道
	service          reposity.UserService // 用户服务
}

// Group 群聊 写入较为频繁 读取基本上都是循环读
type Group struct {
	lock        *sync.Mutex         // 控制群聊
	groupmember map[string]struct{} // 聊天成员
}

// ConnectIdentify 连接情况
type ConnectIdentify struct {
	connect         *net.TCPConn // 连接情况
	expireat        *time.Timer  // 计时器
	ifinconnectpool bool         // 是否加入了连接池
	useraccount     string       // 用户账号
}

// friendMakeInfo 交友信息
type friendMakeInfo struct {
	launcher   int64  //发起者
	accepter   int64  // 接受者
	randomcode string // 随机验证码
}

func NewTcpServer(ctx context.Context) (*TcpServer, error) {

	tcpServerCtx, tcpServerCancel := context.WithCancel(ctx)

	var addr = net.TCPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: Port,
	}

	listener, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		logServer.Error("建立tcp监听失败，失败原因为(%s)", err.Error())
		tcpServerCancel()
		return nil, err
	}
	logServer.Info("成功建立Tcp服务器")
	var service = reposity.UserService{
		ChattingReposity:     xormuser.UserRepository{XormEngine: xormdatabase.DBEngine},
		ChatingCacheReposity: redisuser.UserCacheRepository{RedisEngine: redisdatabase.RedisClient},
	}

	tcpserver := &TcpServer{
		addr:             addr,
		ctx:              tcpServerCtx,
		cancel:           tcpServerCancel,
		listener:         *listener,
		connectionpool:   sync.Map{},
		friendMakeList:   sync.Map{},
		service:          service,
		groupchatting:    sync.Map{},
		groupmessagechan: make(chan Message, 100),
	}

	grouplist, err := tcpserver.service.QueryAllGroup()
	if err != nil {
		return nil, err
	}

	for i := range grouplist {
		var group = Group{
			groupmember: make(map[string]struct{}),
			lock:        &sync.Mutex{},
		}
		groupid := strconv.FormatInt(grouplist[i].Groupid, 10)
		tcpserver.groupchatting.Store(groupid, group)

	}

	go tcpserver.accept()
	go tcpserver.monitor()
	go tcpserver.monitorgroupchat()
	return tcpserver, nil
}

// close tcp关闭时进行的操作
func (tcpserver *TcpServer) close() {
	tcpserver.listener.Close()
	tcpserver.cancel()
	tcpserver.connectionpool.Range(func(key, value interface{}) bool {
		if err := value.(*net.TCPConn).Close(); err != nil {
			logServer.Error("断开连接失败:%s", err.Error())
		}

		tcpserver.connectionpool.Delete(key)
		return true
	})
	logServer.Info("Tcpserver停止服务...")
}

// accept tcp监听程序
func (tcpserver *TcpServer) accept() {
	logServer.Info("开启Tcp监听服务...")
	for {
		connect, err := tcpserver.listener.AcceptTCP()
		if err != nil {
			logServer.Error("连接情况异常:（%s）", err.Error())
			return

		} else {
			// 进行通信 包括转发信息等
			logServer.Info("监听到连接：ip为(%s)", connect.RemoteAddr())
			go tcpserver.chattingWithConnect(connect)
		}

	}

}

// monitor tcp监控程序
func (tcpserver *TcpServer) monitor() {
	defer tcpserver.close()
	logServer.Info("开启Tcp监控服务...")

	<-tcpserver.ctx.Done()
	logServer.Info("tcpserver退出监控服务...")
}

// monitorgroupchat 负责对群聊的管理
func (tcpserver *TcpServer) monitorgroupchat() {
	logServer.Info("群聊监控已开启")
	for {
		select {
		case <-tcpserver.ctx.Done():
			logServer.Info("tcpserver退出监控群聊服务...")
			return
		case message := <-tcpserver.groupmessagechan:
			go tcpserver.dealWithGroupMessage(message)
		}
	}
}

//chattingWithConnect 和具体的连接进行通信
func (tcpserver *TcpServer) chattingWithConnect(connect *net.TCPConn) {

	// 定义是否是已加入连接池
	// var ifaddtoconnctpool = false
	var connectidentify = ConnectIdentify{
		connect:         connect,
		ifinconnectpool: false,
		expireat:        time.NewTimer(30 * time.Second),
		useraccount:     "unlogined",
	}

	// 开启计时器
	go connectidentify.monitorconnect()

	for {

		var receivedData = make([]byte, 1024*2)
		count, err := connectidentify.connect.Read(receivedData)

		if err != nil {
			logServer.Error("接收数据失败：（%s）", err.Error())
			logServer.Info("关闭此连接...")
			if connectidentify.useraccount == "unlogined" {
				connectidentify.connect.Close()
			} else {
				conni, ok := tcpserver.connectionpool.Load(connectidentify.useraccount)
				if ok {
					conn := conni.(*net.TCPConn)
					if conn == connectidentify.connect {
						tcpserver.connectionpool.Delete(connectidentify.useraccount)
						connectidentify.connect.Close()
						useraccountint, _ := strconv.ParseInt(connectidentify.useraccount, 10, 64)
						if err := tcpserver.service.UpdateUserOnlineStatus(useraccountint, false); err != nil {
							logServer.Error("用户修改状态失败:%s", err.Error())
						}
						// 广播下线信息
						userfriend, _ := tcpserver.service.QueryFriends(useraccountint)
						for i := range userfriend.Friends {
							friendaccount := strconv.FormatInt(userfriend.Friends[i].UserAccount, 10)
							conni, ok := tcpserver.connectionpool.Load(friendaccount)
							if ok {
								conn := conni.(*net.TCPConn)
								SendCommonMessage(conn, connectidentify.useraccount, friendaccount, "", "i am offline", NotOlineStatus)
							}
						}
						logServer.Info("用户：（%s）断开连接", connectidentify.useraccount)
					} else {
						logServer.Info("用户:%s重复登录", connectidentify.useraccount)
						connectidentify.connect.Close()
					}
				} else {
					connectidentify.connect.Close()
				}

			}

			return
		}

		var ReceivedStruct Message
		err = json.Unmarshal(receivedData[:count], &ReceivedStruct)

		if err != nil {
			logServer.Error("解析数据失败：（%s）", err.Error())
		} else {
			tcpserver.dealWithMessage(ReceivedStruct, &connectidentify)
		}

	}

}

//dealWithMessage tcp 信息转发 及 信息处理
func (tcpserver *TcpServer) dealWithMessage(receiveMessage Message, connectidentify *ConnectIdentify) {
	var service = reposity.UserService{
		ChattingReposity:     xormuser.UserRepository{XormEngine: xormdatabase.DBEngine},
		ChatingCacheReposity: redisuser.UserCacheRepository{RedisEngine: redisdatabase.RedisClient},
	}

	conn := connectidentify.connect

	useraccount, err := strconv.ParseInt(receiveMessage.Sender, 10, 64)
	if err != nil {
		logServer.Error("用户解析失败: %s", err.Error())
		conn.Close()
		return
	}
	// 鉴权
	ifsame, err := service.IfTokenSameAndNotExpired(useraccount, receiveMessage.Token)

	if err != nil || !ifsame {
		logServer.Error("用户鉴权未通过。")
		SendCommonMessage(conn, "tcpserver provider", receiveMessage.Sender, receiveMessage.Message, "", AuthorizationFail)
		conn.Close()
		return
	}

	// 信息类型
	switch receiveMessage.MessageType {
	case FirstConnect: // 初次连接
		tcpserver.NewUserLoginIn(service, useraccount, receiveMessage, connectidentify)
	case SendMessage: // 发送信息
		tcpserver.SendMessageToReceiver(receiveMessage, conn, SendInfoSuccess, UserNotOnline)
	case SendGroupMessage: // 群聊信息处理
		tcpserver.groupmessagechan <- receiveMessage
	case HeartBeat: // 心跳服务
		tcpserver.HeartBeatMessage(receiveMessage, connectidentify)
	case CloseConnect: // 关闭连接
		tcpserver.CloseConnect(receiveMessage, conn, "已收到断开请求，正在断开...")
	case SendFriendRequest: // 发送好友申请
		tcpserver.LaunchFriendRequest(receiveMessage, conn)
	case AcceptFriendRequest: // 接受好友申请
		tcpserver.AcceptFriendRequest(service, receiveMessage, conn)
	case RejectFriendRequest: // 拒绝好友申请
		tcpserver.RejectFrienRequest(receiveMessage, conn)
	default:
		logServer.Info("未知命令类型")
	}
}

// dealWithGroupMessage 群聊信息的处理
func (tcpserver *TcpServer) dealWithGroupMessage(message Message) {
	thegroupi, ok := tcpserver.groupchatting.Load(message.Groupid)
	if !ok {
		return
	}
	thegroup := thegroupi.(Group)
	mb, err := json.Marshal(message)
	if err != nil {
		logServer.Error("信息转码失败:%s", err.Error())
		return
	}
	thegroup.lock.Lock()
	for k := range thegroup.groupmember {
		if k == message.Sender {
			continue
		}
		conni, ok := tcpserver.connectionpool.Load(k)
		if ok {
			conn := conni.(*net.TCPConn)
			if _, err := conn.Write(mb); err != nil {
				logServer.Error("发送信息至:%v失败:%s", k, err.Error())
			}
		}
	}
	thegroup.lock.Unlock()
	senderi, ok := tcpserver.connectionpool.Load(message.Sender)
	if !ok {
		logServer.Error("发送信息者不存在")
		return
	}
	senderconn := senderi.(*net.TCPConn)
	SendCommonMessage(senderconn, "tcpserver provider", message.Sender, message.Message, message.Groupid, SendGroupInfoSuccess)

}

//SendCommonMessage 用于发送信息
func SendCommonMessage(conn *net.TCPConn, sender, receiver, message, groupid string, receiveStatus MessageTypes) {
	var messageforsend = Message{
		MessageType: receiveStatus,
		Message:     message,
		Token:       "",
		Sender:      sender,
		Receiver:    receiver,
		SendTime:    time.Now(),
		Groupid:     groupid,
	}
	sendbyte, _ := json.Marshal(messageforsend)
	_, err := conn.Write(sendbyte)
	if err != nil {
		logServer.Error("发送回复信息失败：(%s)", err.Error())
	}
}

// monitorconnect 监控连接状态 用于设置超时和垃圾连接
func (ci *ConnectIdentify) monitorconnect() {
	<-ci.expireat.C
	logServer.Info("连接已超时,断开此连接,用户号为:%s,用户登录状态为:%v", ci.useraccount, ci.ifinconnectpool)
	ci.connect.Close()
}
