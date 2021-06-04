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
	addr           net.TCPAddr            // 服务器地址
	conn           map[string]net.TCPConn // 用户连接
	connectionpool sync.Map               // 用户连接池
	listener       net.TCPListener        // tcp监听器
	ctx            context.Context        // 上下文
	cancel         context.CancelFunc     // 退出回调
	friendMakeList []friendMakeInfo
}

// friendMakeInfo 交友信息
type friendMakeInfo struct {
	launcher   int64  //发起者
	accepter   int64  // 接受者
	randomcode string // 随机验证码
}

// connection 用于建立心跳
type connection struct {
	conn net.TCPConn // 连接
	time time.Timer  // 定时器
}

func NewTcpServer(ctx context.Context) {

	tcpServerCtx, tcpServerCancel := context.WithCancel(ctx)

	var addr = net.TCPAddr{
		IP:   net.IPv4(0, 0, 0, 0),
		Port: Port,
	}

	listener, err := net.ListenTCP("tcp", &addr)
	if err != nil {
		logServer.Error("建立tcp监听失败，失败原因为(%s)", err.Error())
		tcpServerCancel()
		return
	}
	logServer.Info("成功建立Tcp服务器")

	tcpserver := &TcpServer{
		addr:           addr,
		ctx:            tcpServerCtx,
		cancel:         tcpServerCancel,
		listener:       *listener,
		conn:           make(map[string]net.TCPConn),
		connectionpool: sync.Map{},
		friendMakeList: make([]friendMakeInfo, 0),
	}
	go tcpserver.accept()
	go tcpserver.monitor()
}

// accept tcp监听程序
func (tcpserver *TcpServer) accept() {
	logServer.Info("开启Tcp监听服务...")
	for {
		select {
		case <-tcpserver.ctx.Done():
			tcpserver.listener.Close()
			logServer.Info("tcpserver停止监听")
			tcpserver.cancel()
			goto stopTcpserverlistener
		default:
			connect, err := tcpserver.listener.AcceptTCP()
			if err != nil {
				logServer.Error("有连接连接至服务器失败（%s）", err.Error())
			} else {
				// 进行通信 包括转发信息等
				logServer.Info("监听到连接：ip为(%s)", connect.RemoteAddr())
				go tcpserver.chattingWithConnect(*connect)
			}

		}

	}
stopTcpserverlistener:
	logServer.Info("tcpserver退出监听服务...")

}

// monitor 监控
func (tcpserver *TcpServer) monitor() {
	logServer.Info("开启Tcp监控服务")
	for {
		select {
		case <-tcpserver.ctx.Done():
			tcpserver.listener.Close()
			logServer.Info("tcpserver停止监控")
			tcpserver.cancel()
			goto stopTcpservermonitor

		}
	}
stopTcpservermonitor:
	logServer.Info("tcpserver退出监控服务...")
}

//chattingWithConnect 和具体的连接进行通信
func (tcpserver *TcpServer) chattingWithConnect(connect net.TCPConn) {
	for {
		var receivedData = make([]byte, 1024*2)
		count, err := connect.Read(receivedData)

		if err != nil {
			logServer.Error("接收数据失败：（%s）", err.Error())
			logServer.Info("关闭此连接...")
			connect.Close()
			return
		}

		var ReceivedStruct Message
		err = json.Unmarshal(receivedData[:count], &ReceivedStruct)

		if err != nil {
			logServer.Error("解析数据失败：（%s）", err.Error())
		} else {
			tcpserver.dealWithMessage(ReceivedStruct, &connect)
		}

	}

}

//dealWithMessage tcp 信息转发 及 信息处理
func (tcpserver *TcpServer) dealWithMessage(receiveMessage Message, conn *net.TCPConn) {
	var service = reposity.UserService{
		ChattingReposity:     xormuser.UserRepository{XormEngine: xormdatabase.DBEngine},
		ChatingCacheReposity: redisuser.UserCacheRepository{RedisEngine: redisdatabase.RedisClient},
	}
	useraccount, err := strconv.ParseInt(receiveMessage.Sender, 10, 64)
	if err != nil {
		logServer.Error("用户解析失败: %s", err.Error())
		conn.Close()
		return
	}
	// 鉴权
	ifsame, err := service.IfTokenSameAndNotExpired(useraccount, receiveMessage.Token)
	if err != nil {
		SendReplyMessage(conn, receiveMessage, AuthorizationFail)
		conn.Close()
		return
	}
	if !ifsame {
		SendReplyMessage(conn, receiveMessage, AuthorizationFail)
		conn.Close()
		return
	}

	// 信息类型
	switch receiveMessage.MessageType {
	case FirstConnect: // 初次连接
		tcpserver.NewUserLoginIn(service, useraccount, receiveMessage, conn)
	case SendMessage: // 发送信息
		tcpserver.SendMessageToReceiver(receiveMessage, conn, SendInfoSuccess, UserNotOnline)
	case HeartBeat: // 心跳服务
		tcpserver.HeartBeatMessage()
	case CloseConnect: // 关闭连接
		tcpserver.CloseConnect(receiveMessage, conn)
	case SendFriendRequest: // 发送好友申请
		tcpserver.LaunchFrienRequest(receiveMessage, conn)
	case AcceptFriendRequest: // 接受好友申请
		tcpserver.AcceptFrienRequest(service, receiveMessage, conn)
	case RejectFriendRequest: // 拒绝好友申请
		tcpserver.RejectFrienRequest(receiveMessage, conn)
	default:
		logServer.Info("未知命令类型")
	}
}

//SendReplyMessage 用于发送 发送状态 的信息
func SendReplyMessage(conn *net.TCPConn, message Message, receiveStatus MessageTypes) {
	var replyMessage = Message{
		MessageType: receiveStatus,
		Message:     message.Message,
		Token:       "",
		Sender:      message.Receiver,
		Receiver:    message.Sender,
		SendTime:    message.SendTime,
	}
	sendbyte, _ := json.Marshal(replyMessage)
	_, err := conn.Write(sendbyte)
	if err != nil {
		logServer.Error("发送回复信息失败：(%s)", err.Error())
	}
}
