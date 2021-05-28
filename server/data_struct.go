//
// Copyright (c) 2021 朱俊杰
// All rights reserved
// filename: dataStruct.go
// description: Tcpserver的数据结构
// version: 0.1.0
// created by zhujunjie(1121883342@qq.com) at 2021-04-27
//

package server

import (
	"context"
	"net"
)

const Port = 4000

// User 用户结构
type User struct {
	Id string // 用户标识
	Conn *net.Conn // 用户连接
	Context context.Context // 上下文用于协程控制
}